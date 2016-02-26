package imager

// #cgo pkg-config: libavcodec libavutil libavformat libswscale
/*

#include <libavcodec/avcodec.h>
#include <libavutil/frame.h>
#include <libavformat/avformat.h>
#include <libswscale/swscale.h>

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
			PIX_FMT_RGB24,
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
				avpicture_alloc((AVPicture *)rgbFrame, PIX_FMT_RGB24, codecCtx->width, codecCtx->height);
				sws_scale(swCtx, frame->data, frame->linesize, 0, frame->height, rgbFrame->data, rgbFrame->linesize);
				rgbFrame->height = frame->height;
				rgbFrame->width = frame->width;
				rgbFrame->format = frame->format;

				AVCodec *outCodec = avcodec_find_encoder(CODEC_ID_PNG);
				AVCodecContext *outCodecCtx = avcodec_alloc_context3(outCodec);
				if (!codecCtx)
					return -1;

				outCodecCtx->width = codecCtx->width;
				outCodecCtx->height = codecCtx->height;
				outCodecCtx->pix_fmt = PIX_FMT_RGBA;
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
*/
import "C"
import "unsafe"

func thumbnailWebm(input, output string) {
	cin := C.CString(input)
	cout := C.CString(output)
	C.convert_first_frame_to_png(cin, cout)
	defer C.free(unsafe.Pointer(cin))
	defer C.free(unsafe.Pointer(cout))
}
