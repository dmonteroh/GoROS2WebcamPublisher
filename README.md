# GoROS2WebcamPublisher
Simple implementation to capture Webcam video and publish it to a ROS2 topic using GOCV

## Requires 
ROS2 Wrapper: github.com/TIERS/rclgo, ROS2 Messages: github.com/TIERS/rclgo-msgs and GO CV wrapper gocv.io/x/gocv

## What is it for?
Made for testing ROS2 nodes, the code takes the computer's webcam image and publishes it to a ROS2 topic. The default topic is /webcam/image, but it can be changed by using CLI flags at run time.

## Usage
Build or Run, and use CLI flags as required

Run:
go run webcam_publisher.go -topic=/test/topic

Build:
./webcamPublisher -topic=/test/topic

## Flags
  -bigendian uint (Not clear what this data is used for, it was added in case a ROS2 implementation required a specific value)
    	Bigendian for the ROS MSG (default 1)
  -device int (Usually the camera device from a laptop works with 0, which is the default number)
    	Webcam Device identification (default 0)
  -namespace string (Can be changed if a specific Namespace for the Node must be used)
    	Namespace for the ROS node (default "tools")
  -node string (Can be changed if a specific Node name is required, like when running multiple instances)
    	ROS node name (default "webcam_to_topic")
  -topic string (Allows for the topic to be changed, no need to change any code)
    	Topic to which the image will be published to (default "/webcam/image")
  -window (CV can create a window to show the image directly from the webcam, helps with testing)
    	Whether or not to show the image captured by the webcam (default false)
