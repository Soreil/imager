package webm

// #cgo pkg-config: libavcodec libavutil libavformat libswscale
/*

#include <libavcodec/avcodec.h>
#include <libavutil/frame.h>
#include <libavformat/avformat.h>
#include <libswscale/swscale.h>

#define PIX_FMT_CHOSEN PIX_FMT_RGBA
#define BUFFER_SIZE 8092

#define CHECK_ERR(ERR) {if ((ERR)<0) return ERR; }

static int readFunction(void* opaque,  unsigned char* buf, int buf_size) {
}

AVFrame * extract_webm_image(unsigned char *opaque,size_t len)
{
	av_register_all();
	avcodec_register_all();

	//TODO(sjon): Buffer memory has to be allocated by av_malloc()
	//free with av_free()
	unsigned char *buffer = (unsigned char*)av_malloc(BUFFER_SIZE);


	//Allocate avioContext
	//Has to be av_free()'d
	//TODO(sjon): Implement custom reader
	AVIOContext *ioCtx = avio_alloc_context(buffer,BUFFER_SIZE,0,opaque,NULL,NULL,NULL);

	//destroy with avformat_free_contex()
	AVFormatContext * ctx = avformat_alloc_context();

	//Set up context to read from memory
	ctx->pb = ioCtx;

	//TODO(sjon):Actually implement the not file based option
	//open takes a fake filename when the context pb field is set up
	int err = avformat_open_input(&ctx, "Lmao, using actual files", NULL, NULL);
	CHECK_ERR(err);

	err = avformat_find_stream_info(ctx,NULL);
	CHECK_ERR(err);

	AVCodec * codec = NULL;
	int strm = av_find_best_stream(ctx, AVMEDIA_TYPE_VIDEO, -1, -1, &codec, 0);

	AVCodecContext * codecCtx = ctx->streams[strm]->codec;
	err = avcodec_open2(codecCtx, codec, NULL);
	CHECK_ERR(err);

	struct SwsContext * swCtx = sws_getContext(codecCtx->width,
			codecCtx->height,
			codecCtx->pix_fmt,
			codecCtx->width,
			codecCtx->height,
			PIX_FMT_CHOSEN,
			SWS_FAST_BILINEAR, NULL, NULL, NULL);

	for (;;)
	{
		AVPacket pkt;
		err = av_read_frame(ctx, &pkt);
		CHECK_ERR(err);

		if (pkt.stream_index == strm)
		{
			int got = 0;
			AVFrame * frame = av_frame_alloc();
			err = avcodec_decode_video2(codecCtx, frame, &got, &pkt);
			CHECK_ERR(err);

			if (got)
			{
				AVFrame * rgbFrame = av_frame_alloc();
				avpicture_alloc((AVPicture *)rgbFrame, PIX_FMT_CHOSEN, codecCtx->width, codecCtx->height);
				sws_scale(swCtx, frame->data, frame->linesize, 0, frame->height, rgbFrame->data, rgbFrame->linesize);
				//TODO(sjon): sws_scale performed, we can dump the pixels at this point to work with as raw RGBA data
				rgbFrame->height = frame->height;
				rgbFrame->width = frame->width;
				rgbFrame->format = frame->format;

				return rgbFrame;
			}
			//av_frame_free(&frame);
		}
	}
}

AVFormatContext * extract_webm_ctx(unsigned char *opaque,size_t len)
{
	av_register_all();
	avcodec_register_all();

	//TODO(sjon): Buffer memory has to be allocated by av_malloc()
	//free with av_free()
	unsigned char *buffer = (unsigned char*)av_malloc(BUFFER_SIZE);


	//Allocate avioContext
	//Has to be av_free()'d
	//TODO(sjon): Implement custom reader
	AVIOContext *ioCtx = avio_alloc_context(buffer,BUFFER_SIZE,0,opaque,NULL,NULL,NULL);

	//destroy with avformat_free_contex()
	AVFormatContext * ctx = avformat_alloc_context();

	//Set up context to read from memory
	ctx->pb = ioCtx;

	//TODO(sjon):Actually implement the not file based option
	//open takes a fake filename when the context pb field is set up
	int err = avformat_open_input(&ctx, "Lmao, using actual files", NULL, NULL);
	CHECK_ERR(err);

	err = avformat_find_stream_info(ctx,NULL);
	CHECK_ERR(err);
	return ctx;
}
*/
import "C"
import (
	"errors"
	"fmt"
	"image"
	"io"
	"unsafe"
)

//TODO(sjon): add actual header here
const webmHeader = ""

func init() {
	image.RegisterFormat("webm", webmHeader, Decode, DecodeConfig)
}

type avFrame struct {
	frame *C.AVFrame
}

//TODO(sjon): Blows up on printing
func decode(data []byte) {
	//frame := C.extract_webm_image((*C.uchar)(unsafe.Pointer(&data[0])), C.size_t(len(data)))
	////fmt.Printf("%+v\n", frame)
	//avFrame := avFrame{frame}
	//fmt.Println(avFrame)
	fmt.Printf("%+v\n", C.extract_webm_ctx((*C.uchar)(unsafe.Pointer(&data[0])), C.size_t(len(data))))
}

//TODO(sjon):Use C code to decode, need to find a way to create a formatcontext without reading file from disk
func Decode(r io.Reader) (image.Image, error) {
	return nil, errors.New("Not implemented")
}

//TODO(sjon):Use C code first part, return before sws_scale
func DecodeConfig(r io.Reader) (image.Config, error) {
	return image.Config{}, errors.New("Not implemented")
}
