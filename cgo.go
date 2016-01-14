package imager

// #cgo pkg-config: libavcodec libavutil
//
// #include "libavcodec/avcodec.h"
// #include "libavutil/frame.h"
import "C"
import (
	"io/ioutil"
	"unsafe"
)

func init() {
	//Register all the codecs, parsers and bitstream filters which were enabled at configuration time.
	//If you do not call this function you can select exactly which formats you want to support, by using the individual registration functions.
	C.avcodec_register_all()
}

type Decoder struct {
	codec    *C.AVCodec
	codecCtx *C.AVCodecContext
}

//VP8 only, need to figure the logic behind the header out, maybe this is a versioning thing?
func NewVP8Decoder() *Decoder {
	var dec Decoder
	dec.codec = C.avcodec_find_decoder(C.AV_CODEC_ID_VP8)
	//if dec.codec == nil {
	//	panic(errors.New("Can't find codec"))
	//}
	dec.codecCtx = C.avcodec_alloc_context3(dec.codec)
	//if dec.codecCtx == nil {
	//	panic(errors.New("Can't create context"))
	//}
	dec.codecCtx.opaque = unsafe.Pointer(&dec)
	C.avcodec_open2(dec.codecCtx, dec.codec, nil)
	return &dec
}

func (d *Decoder) DecodeVP8(fileName string) error {
	var packet *C.AVPacket
	var frame *C.AVFrame
	var got C.int

	frame = C.av_frame_alloc()
	defer C.av_frame_free(&frame)

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	packet = (*C.AVPacket)(C.malloc(C.size_t(unsafe.Sizeof(C.AVPacket{}))))
	defer C.free(unsafe.Pointer(packet))
	C.av_init_packet(packet)

	packet.data = (*C.uint8_t)(&file[0])
	packet.size = C.int(len(file))

	C.av_read_frame(s, packet)
	//if C.avcodec_decode_video2(d.codecCtx, frame, &got, packet) < 0 {
	//	return errors.New("Unable to decode")
	//}
	//if got == 0 {
	//	return errors.New("Didn't get any data")
	//}
	return nil
}
