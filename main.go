package imager

//type media interface{}
//type mediaType int
//
//const (
//	png mediaType = iota
//	jpg
//	gif
//	webm
//)
//
//func loadMedia(s string) (media error) {
//	med, err := open(s)
//	if err != nil {
//		panic(err)
//	}
//	isValid(med)
//	medType := typeOf(med)
//
//	processImage(med, medType)
//}
//
//func processMedia(med, medType) {
//	switch medType {
//	case png:
//		pngProcessor <- med
//	case jpg:
//		jpgProcessor <- med
//	case gif:
//		gifProcessor <- med
//	}
//}
//
//func scheduler() {
//	for {
//		select {
//		case <-aWorkerIsDone:
//			addworkToIt()
//		default:
//			stall()
//		}
//	}
//}
