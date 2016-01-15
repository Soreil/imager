package imager

// #cgo pkg-config: libavcodec libavutil libavformat
//
// #include "libavcodec/avcodec.h"
// #include "libavutil/frame.h"
// #include "libavformat/avformat.h"
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
	var pict *C.AVPicture = nil
	var got C.int

	cs := C.CString(fileName)
	defer C.free(unsafe.Pointer(cs))

	if r := C.avformat_open_input(&formatCtx, cs, nil, nil); r < 0 {
		return errors.New("Failed to open input file" + fileName)
	}

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

	if r := C.avpicture_alloc(pict, d.codecCtx.pix_fmt, frame.width, frame.height); r < 0 {
		return errors.New("Failed to allocate picture buffer")
	}
	defer C.avpicture_free(pict)
	if r := C.avpicture_fill(pict, frame.data, d.codecCtx.pix_fmt, frame.width, frame.height); r < 0 {
		return errors.New("Failed to fill picture")
	}

	return nil
}
