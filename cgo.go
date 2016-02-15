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
	codec     *C.AVCodec
	codecCtx  *C.AVCodecContext
	formatCtx *C.AVFormatContext
}

//VP8 only, need to figure the logic behind the header out, maybe this is a versioning thing?
func NewVP8Decoder() *Decoder {
	//var dec Decoder
	//dec.codec = C.avcodec_find_decoder(C.AV_CODEC_ID_VP8)
	//if dec.codec == nil {
	//	panic(errors.New("Can't find codec"))
	//}
	//dec.codecCtx = C.avcodec_alloc_context3(dec.codec)
	//if dec.codecCtx == nil {
	//	panic(errors.New("Can't create context"))
	//}
	//dec.codecCtx.opaque = unsafe.Pointer(&dec)
	//C.avcodec_open2(dec.codecCtx, dec.codec, nil)
	//return &dec
	return &Decoder{}
}

func (d *Decoder) DecodeVP8(fileName string) error {

	cs := C.CString(fileName)
	defer C.free(unsafe.Pointer(cs))

	if r := C.avformat_open_input(&d.formatCtx, cs, nil, nil); r < 0 {
		return errors.New("Failed to open input file" + fileName)
	}
	if r := C.avformat_find_stream_info(d.formatCtx, nil); r < 0 {
		return errors.New("Stream info broken")
	}
	if r := C.av_find_best_stream(d.formatCtx, C.AVMEDIA_TYPE_VIDEO, -1, -1, &d.codec, 0); r < 0 {
		return errors.New("Failed to locate best video codec")
	}

	//strm is an offset for the d.formatCtx array, assuming it is always 0
	//in our use cases skipping it should be fine right?

	strm := C.av_find_best_stream(d.formatCtx, C.AVMEDIA_TYPE_VIDEO, -1, -1, &d.codec, 0)
	if strm != 0 {
		return errors.New("Video stream not at position zero:" + fmt.Sprint(strm))
	}

	d.codecCtx = (**d.formatCtx.streams).codec
	//fmt.Printf("%+v\n", d.codecCtx)

	//THIS CONFLICTS WITH THE DECODER NEW FUNCTION
	if r := C.avcodec_open2(d.codecCtx, d.codec, nil); r < 0 {
		return errors.New("Failed to open codec")
	}
	// ... or does it? At least this here too "works"
	//fmt.Printf("%+v\n", d.codecCtx)

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

	//DECODING PART OF THE FUNCTION

	//Allocate frame data receiver packet
	var packet C.AVPacket

	if r := C.av_read_frame(d.formatCtx, &packet); r < 0 {
		return errors.New("Failed to initialize packet from formatContext")
	}

	//Dead pointer code
	//packet = (*C.AVPacket)(C.malloc(C.size_t(unsafe.Sizeof(C.AVPacket{}))))
	//defer C.free(unsafe.Pointer(packet))
	//C.av_init_packet(packet)

	//Allocate the frame itself
	var frame *C.AVFrame = C.av_frame_alloc()
	defer C.av_frame_free(&frame)

	var got C.int
	if C.avcodec_decode_video2(d.codecCtx, frame, &got, &packet) < 0 {
		return errors.New("Unable to decode")
	} else if got == 0 {
		return errors.New("Didn't get any data from decoding video")
	}

	fmt.Println("Video width:", frame.width)
	fmt.Println("Video height:", frame.height)

	//Allocate output frame
	var pngFrame *C.AVFrame = C.av_frame_alloc()
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
