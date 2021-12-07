package main

import (
	"flag"
	"log"
	"sync"

	sensor_msgs "github.com/TIERS/rclgo-msgs/sensor_msgs/msg"
	"github.com/TIERS/rclgo/pkg/rclgo"
	"gocv.io/x/gocv"
)

var imageTopic string
var namespace string
var nodeName string
var showWindow bool
var captureDevice int
var bigendian uint8 = 1
var window *gocv.Window

func main() {
	//CLI flags to configure the script, and default values
	topicPtr := flag.String("topic", "/webcam/image", "Topic to which the image will be published to")
	namespacePtr := flag.String("namespace", "tools", "Namespace for the ROS node")
	namePtr := flag.String("node", "webcam_to_topic", "ROS node name")
	windowPtr := flag.Bool("window", false, "Whether or not to show the image captured by the webcam")
	devicePtr := flag.Int("device", 0, "Webcam Device identification")
	bigendianPtr := flag.Uint("bigendian", 1, "Bigendian for the ROS MSG")
	flag.Parse()

	//Setting flags to variables
	if topicPtr != nil && *topicPtr != imageTopic && *topicPtr != "" {
		imageTopic = *topicPtr
	}

	if namespacePtr != nil && *namespacePtr != namespace && *namespacePtr != "" {
		namespace = *namespacePtr
	}

	if namePtr != nil && *namePtr != nodeName && *namePtr != "" {
		nodeName = *namePtr
	}

	if windowPtr != nil && *windowPtr != showWindow {
		showWindow = *windowPtr
	}

	if devicePtr != nil && *devicePtr != captureDevice && *devicePtr >= 0 {
		captureDevice = *devicePtr
	}

	if bigendianPtr != nil && *bigendianPtr != uint(bigendian) {
		bigendian = uint8(*bigendianPtr)
	}

	//Starting the Webcam video capture
	webcam, err := gocv.VideoCaptureDevice(captureDevice)
	if err != nil {
		log.Fatalf("Unable to initialize Webcam: %v", err)
	}
	defer webcam.Close()

	//Create the variable where CV will store the image
	img := gocv.NewMat()
	defer img.Close()

	//If the Image is going to be shown, initialize the window variable
	if showWindow {
		window = gocv.NewWindow("Webcam Capture")
	}
	defer window.Close()

	//Setting up ROS Publisher
	var doneChannel = make(chan bool)
	var wg sync.WaitGroup
	rclArgs, rclErr := rclgo.NewRCLArgs("")
	if rclErr != nil {
		log.Fatal(rclErr)
	}

	rclContext, rclErr := rclgo.NewContext(&wg, 0, rclArgs)
	if rclErr != nil {
		log.Fatal(rclErr)
	}
	defer rclContext.Close()

	rclNode, rclErr := rclContext.NewNode(nodeName, namespace)
	if rclErr != nil {
		log.Fatal(rclErr)
	}

	opts := rclgo.NewDefaultPublisherOptions()
	opts.Qos.Reliability = rclgo.RmwQosReliabilityPolicySystemDefault
	pub, err := rclNode.NewPublisher(imageTopic, sensor_msgs.ImageTypeSupport, opts)
	if err != nil {
		log.Fatalf("Unable to create publisher: %v", err)
	}

	//Infinite loop
	for {
		//If the webcam can be read, then update img variable
		if ok := webcam.Read(&img); !ok || img.Empty() {
			log.Print("Unable to read from webcam")
			continue
		}
		//Call the image bridge function
		rosmsg := imageBridge(img, "passthrough", bigendian)
		pub.Publish(rosmsg)

		//Updates the Window image if necessary
		if showWindow {
			window.IMShow(img)
			window.WaitKey(1)
		}
	}
	<-doneChannel
}

//This function takes the CV image Mat and returns a ROS Image Message
func imageBridge(image gocv.Mat, encoding string, bigendian uint8) *sensor_msgs.Image {
	img_msg := sensor_msgs.NewImage()
	img_msg.Height = uint32(image.Rows())
	img_msg.Width = uint32(image.Cols())
	img_msg.Encoding = encoding
	//I have no idea what this is, but it's worked with numbers 0 to 2 just fine
	img_msg.IsBigendian = bigendian
	img_msg.Data = image.ToBytes()
	img_msg.Step = uint32(len(img_msg.Data))
	return img_msg
}
