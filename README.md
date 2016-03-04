# imager
High performance thumbnailer for integration with Go programs.

##Dependencies:
 -WebM driver
  +libavcodec 
  +libavutil 
  +libavformat 
  +libswscale
 -SVG driver
  +rsvg-convert
 -PDF driver
  +ghostscript built with PDF support

 -PNG compressor
  +pngquant

##Usage:
The package can be called with an io.Reader containing a valid  supported media type and a desired thumbnail size. Depending on the features of the filetype given the output will be PNG (transparency) or JPEG.
In case of a desire to use other drivers for the formats supported than those maintained by me you will have different dependencies for those parts.

Manually setting input to output file type associations and compression levels (per file) will possibly be implemented in the future. 
