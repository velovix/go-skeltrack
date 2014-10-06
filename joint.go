package skeltrack

/*
#include <skeltrack.h>
*/
import "C"

// JointID represents a joint.
type JointID int

// All available joint values.
const (
	JointHead          = JointID(C.SKELTRACK_JOINT_ID_HEAD)
	JointLeftShoulder  = JointID(C.SKELTRACK_JOINT_ID_LEFT_SHOULDER)
	JointRightShoulder = JointID(C.SKELTRACK_JOINT_ID_RIGHT_SHOULDER)
	JointLeftElbow     = JointID(C.SKELTRACK_JOINT_ID_LEFT_ELBOW)
	JointRightElbow    = JointID(C.SKELTRACK_JOINT_ID_RIGHT_ELBOW)
	JointLeftHand      = JointID(C.SKELTRACK_JOINT_ID_LEFT_HAND)
	JointRightHand     = JointID(C.SKELTRACK_JOINT_ID_RIGHT_HAND)
)

// String returns a string representation of a JointID constant.
func (jointID JointID) String() string {

	switch jointID {
	case JointHead:
		return "head"
	case JointLeftShoulder:
		return "left shoulder"
	case JointRightShoulder:
		return "right shoulder"
	case JointLeftElbow:
		return "left elbow"
	case JointRightElbow:
		return "right elbow"
	case JointLeftHand:
		return "left hand"
	case JointRightHand:
		return "right hand"
	default:
		return "invalid"
	}
}

// Joint is an object that represents a successfully tracked human joint. Joint
// objects should not be created outside of the library.
type Joint struct {
	joint *C.SkeltrackJoint
}

// Type returns the type of joint the Joint object represents.
func (joint *Joint) Type() JointID {

	if joint.joint == nil {
		panic("cannot get the type of a nil Joint object")
	}

	return JointID(joint.joint.id)
}

// Coords returns the real-world coordinates of the Joint object in
// millimeters.
func (joint *Joint) Coords() (int, int, int) {

	if joint.joint == nil {
		panic("cannot get the coordinates of a nil Joint object")
	}

	return int(joint.joint.x), int(joint.joint.y), int(joint.joint.z)
}

// ScreenCoords returns the screen coordinates of the Joint object to make
// creating a skeleton display on-screen easier.
func (joint *Joint) ScreenCoords() (int, int) {

	if joint.joint == nil {
		panic("cannot get the screen coordinates of a nil Joint object")
	}

	return int(joint.joint.screen_x), int(joint.joint.screen_y)
}

// Free frees the memory of the Joint object. The Joint object is unusable
// after this point. This should be done before a Joint object falls out of
// scope.
func (joint *Joint) Free() {

	if joint.joint != nil {
		C.skeltrack_joint_free(joint.joint)
	}
}
