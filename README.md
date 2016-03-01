# imager
CGo FFmpeg bindings for thumbnailing WebM video.
CGo PNGQuant bindings for reducing PNG file sizes.

The goal is to have GIF,PNG,WEBM,JPEG all be represented by a single intermediary pixel format and raw pixel data.
This will be worked on by the Go image package, the idea for now it to have it all be RGBA.

This will be scaled to 250x250 for low quality thumb and 500x500 for high quality.

Only JPEG has no transparency, for that the output format is JPEG too because of that
The rest all need PNG thumbnails to maintain transparency.
Well GIF could have GIF thumbnails.

