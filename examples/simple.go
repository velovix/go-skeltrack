// This is a simple example of a program that uses go-skeltrack. The program
// initializes a Kinect using a binding of the libfreenect library, processes
// the depth information through go-skeltrack, and lets the user know the first
// time a joint is found. This demo assumes that you're using a Kinect, but
// the skeltrack library is device-agnostic.
package main

import (
	"fmt"
	"github.com/velovix/go-freenect"
	"github.com/velovix/go-skeltrack"
	"os"
	"time"
)

const (
	ThresholdBegin = 50
	ThresholdEnd   = 1500
)

var (
	context     freenect.Context
	kinect      freenect.Device
	skeleton    skeltrack.Skeleton
	foundJoints map[skeltrack.JointID]bool
	frameCnt    int
	stop        bool
)

func scaleDepth(data []uint16, width, height, scaledValue int) []uint16 {

	scaledWidth := (width - width%scaledValue) / scaledValue
	scaledHeight := (height - height%scaledValue) / scaledValue

	scaledData := make([]uint16, scaledWidth*scaledHeight)

	for x := 0; x < scaledWidth; x++ {
		for y := 0; y < scaledHeight; y++ {
			index := y*width*scaledValue + x*scaledValue
			value := data[index]

			// Gets rid of any depth values outside of the threshold constants. This
			// helps remove any noise created by unnecessary background depth info.
			if value < ThresholdBegin || value > ThresholdEnd {
				scaledData[y*scaledWidth+x] = 0
			} else {
				scaledData[y*scaledWidth+x] = value
			}
		}
	}

	return scaledData
}

func onDepthFrame(device *freenect.Device, depth []uint16, timestamp uint32) {

	frameCnt++

	// Scale the depth to a more managable size. For the Kinect, the default depth
	// resolution is 640x480, which is too big to process quickly. 16 is the default
	// that go-skeltrack expects, so be sure to notify the library if you decide to
	// scale the depth data differently by using skeleton.SetDimensionReduction.
	depthScaled := scaleDepth(depth, 640, 480, 16)

	// This is where skeleton tracking takes place. The returned value is a map
	// that contains all the joints that were found. Specific joints can be retrieved
	// using JointID constants. i.e. list[skeltrack.JointLeftHand]
	list, err := skeleton.TrackJoints(depthScaled, 640/16, 480/16)
	if err != nil {
		fmt.Println(len(depth))
		fmt.Println(err)
		os.Exit(1)
	}

	// Look to see if any new joints have been found
	for _, val := range list {
		if !foundJoints[val.Type()] {
			fmt.Println("Found joint", val.Type())
			foundJoints[val.Type()] = true
		}

		// Clean up the Joint object
		val.Free()
	}

	// Check if all joints have been found yet, and exit if so
	if len(foundJoints) == 7 {
		fmt.Println("Found all joints!")
		stop = true
	}
}

func init() {

	var err error

	// Create a freenect context.
	context, err = freenect.NewContext()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Set the log level to be fairly verbose.
	context.SetLogLevel(freenect.LogDebug)

	// Open the kinect device
	if cnt, _ := context.DeviceCount(); cnt == 0 {
		fmt.Println("could not find any devices")
		os.Exit(1)
	}
	kinect, err = context.OpenDevice(0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Set the depth callback. This function is called whenever a depth frame is retrieved
	kinect.SetDepthCallback(onDepthFrame)

	// Start the depth camera and begin streaming
	err = kinect.StartDepthStream(freenect.ResolutionMedium, freenect.DepthFormatMM)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initializes the Skeleton object. This is necessary before use!
	skeleton = skeltrack.NewSkeleton()

	foundJoints = make(map[skeltrack.JointID]bool)
}

func main() {

	initTime := time.Now()

	// Waits for new joints until 10 seconds have passed or all joints have been found
	for time.Since(initTime).Seconds() < 10.0 && !stop {
		// Process freenect events
		err := context.ProcessEvents(0)
		if err != nil {
			fmt.Println(err)
			break
		}
	}

	fmt.Println("Processed", frameCnt, "frames in", time.Since(initTime).Seconds(), "seconds.")

	// Clean up objects
	kinect.StopDepthStream()
	kinect.Destroy()
	context.Destroy()
}
