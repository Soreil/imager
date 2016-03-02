//Converts media types in to thumbnails
package imager

// #cgo pkg-config: libavcodec libavutil libavformat libswscale
// #cgo CFLAGS: -Wno-int-conversion -Wno-incompatible-pointer-types
/*
#include <libavcodec/avcodec.h>
#include <libavutil/frame.h>
#include <libavformat/avformat.h>
#include <libswscale/swscale.h>

#define PIX_FMT_CHOSEN PIX_FMT_RGBA

#define CHECK_ERR(ERR) {if ((ERR)<0) return -1; }

int convert_first_frame_to_png(const char *inputVideoFileName, const char *outputPngName)
{
	av_register_all();
	avcodec_register_all();

	AVFormatContext * ctx = NULL;
	int err = avformat_open_input(&ctx, inputVideoFileName, NULL, NULL);
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

				AVCodec *outCodec = avcodec_find_encoder(CODEC_ID_PNG);
				AVCodecContext *outCodecCtx = avcodec_alloc_context3(outCodec);
				if (!codecCtx)
					return -1;

				outCodecCtx->width = codecCtx->width;
				outCodecCtx->height = codecCtx->height;
				outCodecCtx->pix_fmt = PIX_FMT_CHOSEN;
				outCodecCtx->codec_type = AVMEDIA_TYPE_VIDEO;
				outCodecCtx->time_base.num = codecCtx->time_base.num;
				outCodecCtx->time_base.den = codecCtx->time_base.den;

				if (!outCodec || avcodec_open2(outCodecCtx, outCodec, NULL) < 0) {
					return -1;
				}

				AVPacket outPacket;
				av_init_packet(&outPacket);
				outPacket.size = 0;
				outPacket.data = NULL;
				int gotFrame = 0;
				int ret = avcodec_encode_video2(outCodecCtx, &outPacket, rgbFrame, &gotFrame);
				if (ret >= 0 && gotFrame)
				{
					FILE * outPng = fopen(outputPngName, "wb");
					fwrite(outPacket.data, outPacket.size, 1, outPng);
					fclose(outPng);
				}

				avcodec_close(outCodecCtx);
				av_free(outCodecCtx);

				break;
			}
			av_frame_free(&frame);
		}
	}
}

AVFrame * convert_first_frame_to_raw(const char *inputVideoFileName)
{
	av_register_all();
	avcodec_register_all();

	AVFormatContext * ctx = NULL;
	int err = avformat_open_input(&ctx, inputVideoFileName, NULL, NULL);
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
*/
import "C"
import (
	"image"
	"unsafe"

	"github.com/nfnt/resize"
)

//Regular quality preset
var normal image.Point = image.Point{X: 250, Y: 250}

//High quality preset
var sharp image.Point = image.Point{X: 500, Y: 500}

//Extracts a frame and writes to file as PNG
func videoToPNG(input, output string) {
	cin := C.CString(input)
	cout := C.CString(output)
	C.convert_first_frame_to_png(cin, cout)
	defer C.free(unsafe.Pointer(cin))
	defer C.free(unsafe.Pointer(cout))
}

////Wrapper around FFmpeg AVFrame
//type avFrame struct {
//	frame *C.AVFrame
//}

//Takes in a file and returns an image
func extractVideoFrame(input string) image.Image {
	cin := C.CString(input)
	defer C.free(unsafe.Pointer(cin))

	f := C.convert_first_frame_to_raw(cin)
	bs := C.GoBytes(unsafe.Pointer(f.data[0]), f.linesize[0]*f.height)
	return &image.RGBA{Pix: bs,
		Stride: int(f.linesize[0]),
		Rect:   image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: int(f.width), Y: int(f.height)}}}
}

//TODO(sjon): evaluate best resizing algorithm
//Resizes the image to max dimensions
func scale(img image.Image, p image.Point) image.Image {
	return resize.Thumbnail(uint(p.X), uint(p.Y), img, resize.Bilinear)
}
