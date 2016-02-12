package imager

// #cgo pkg-config: libavcodec libavutil libavformat libswscale
//
// #include "libavcodec/avcodec.h"
// #include "libavutil/frame.h"
// #include "libavformat/avformat.h"
// #include "libswscale/swscale.h"
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

func init() {
	//Register all the codecs, parsers and bitstream filters which were enabled at configuration time.
	//If you do not call this function you can select exactly which formats you want to support, by using the individual registration functions.
	C.avcodec_register_all()
	C.av_register_all()
}

type Decoder struct {
	codec    *C.AVCodec
	codecCtx *C.AVCodecContext
}

//VP8 only, need to figure the logic behind the header out, maybe this is a versioning thing?
func NewVP8Decoder() *Decoder {
	var dec Decoder
	dec.codec = C.avcodec_find_decoder(C.AV_CODEC_ID_VP8)
	if dec.codec == nil {
		panic(errors.New("Can't find codec"))
	}
	dec.codecCtx = C.avcodec_alloc_context3(dec.codec)
	if dec.codecCtx == nil {
		panic(errors.New("Can't create context"))
	}
	dec.codecCtx.opaque = unsafe.Pointer(&dec)
	C.avcodec_open2(dec.codecCtx, dec.codec, nil)
	return &dec
}

func (d *Decoder) DecodeVP8(fileName string) error {
	var packet *C.AVPacket
	var frame *C.AVFrame
	var formatCtx *C.AVFormatContext = nil
	var got C.int

	cs := C.CString(fileName)
	defer C.free(unsafe.Pointer(cs))

	if r := C.avformat_open_input(&formatCtx, cs, nil, nil); r < 0 {
		return errors.New("Failed to open input file" + fileName)
	}
	if r := C.av_find_best_stream(formatCtx, C.AVMEDIA_TYPE_VIDEO, -1, -1, &d.codec, 0); r < 0 {
		return errors.New("Failed to locate best video codec")
	}

	//strm := C.av_find_best_stream(formatCtx, C.AVMEDIA_TYPE_VIDEO, -1, -1, &d.codec, 0)
	fmt.Printf("%+v\n", d.codecCtx)
	d.codecCtx = (**formatCtx.streams).codec
	fmt.Printf("%+v\n", d.codecCtx)

	//SWSCONTEXT
	var swCtx *C.struct_SwsContext
	swCtx = C.sws_getContext(d.codecCtx.width,
		d.codecCtx.height,
		d.codecCtx.pix_fmt,
		d.codecCtx.width,
		d.codecCtx.height,
		C.PIX_FMT_RGB24,
		C.SWS_FAST_BILINEAR, nil, nil, nil)
	//SWSCONTEXT

	frame = C.av_frame_alloc()
	defer C.av_frame_free(&frame)

	packet = (*C.AVPacket)(C.malloc(C.size_t(unsafe.Sizeof(C.AVPacket{}))))
	defer C.free(unsafe.Pointer(packet))
	C.av_init_packet(packet)

	C.av_read_frame(formatCtx, packet)
	if C.avcodec_decode_video2(d.codecCtx, frame, &got, packet) < 0 {
		return errors.New("Unable to decode")
	}
	if got == 0 {
		return errors.New("Didn't get any data")
	}
	fmt.Println("Video width:", frame.width)
	fmt.Println("Video height:", frame.height)

	//Allocate output frame
	var pngFrame *C.AVFrame = C.avcodec_alloc_frame()
	//Allocate picture sub thing of output frame
	C.avpicture_alloc((*C.AVPicture)(unsafe.Pointer(&pngFrame)), C.PIX_FMT_RGB24, d.codecCtx.width, d.codecCtx.height)
	//sws_scale in to the new frame
	C.sws_scale(swCtx, &frame.data[0], &frame.linesize[0], 0, frame.height, &pngFrame.data[0], &pngFrame.linesize[0])
	//Find PNG encoder
	var outCodec *C.AVCodec = C.avcodec_find_decoder(C.CODEC_ID_PNG)
	var outCodecCtx *C.AVCodecContext = C.avcodec_alloc_context3(d.codec)
	if outCodecCtx == nil {
		return errors.New("Failed to allocate output codec")
	}
	//Set out output context
	outCodecCtx.width = d.codecCtx.width
	outCodecCtx.height = d.codecCtx.height
	outCodecCtx.pix_fmt = C.PIX_FMT_RGB24
	outCodecCtx.codec_type = C.AVMEDIA_TYPE_VIDEO
	outCodecCtx.time_base.num = d.codecCtx.time_base.num
	outCodecCtx.time_base.den = d.codecCtx.time_base.den
	//open encoder
	if outCodec == nil || C.avcodec_open2(outCodecCtx, outCodec, nil) < 0 {
		return errors.New("Failed to open codec")
	}
	//set up output packet
	var outPacket *C.AVPacket
	outPacket = (*C.AVPacket)(C.malloc(C.size_t(unsafe.Sizeof(C.AVPacket{}))))
	defer C.free(unsafe.Pointer(outPacket))
	C.av_init_packet(outPacket)
	//encode png in to output packet
	got = 0
	if C.avcodec_encode_video2(outCodecCtx, outPacket, pngFrame, &got) < 0 || got == 0 {
		return errors.New("Failed to encode PNG")
	}
	//Write out output packet to disk
	fmt.Println(outPacket.data)
	//EXTERNAL
	//Create cropped to thumbnail PNG
	//Compress PNG using pngquant (as lib?)

	return nil
}
