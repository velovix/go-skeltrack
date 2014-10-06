package skeltrack

/*
#include <skeltrack.h>
#include <stdbool.h>

getObjParamUint(gpointer object, const char *name, uint *val) {

	g_object_get(object, name, val, NULL);
}

getObjParamBool(gpointer object, const char *name, bool *val) {

	g_object_get(object, name, val, NULL);
}

getObjParamFloat(gpointer object, const char *name, float *val) {

	g_object_get(object, name, val, NULL);
}

setObjParamUint(gpointer object, const char *name, uint val) {

	g_object_set(object, name, val, NULL);
}

setObjParamBool(gpointer object, const char *name, bool val) {

	g_object_set(object, name, val, NULL);
}

setObjParamFloat(gpointer object, const char *name, float val) {

	g_object_set(object, name, val, NULL);
}
*/
import "C"
import "errors"
import "unsafe"
import "reflect"

// Skeleton contains and handles skeleton tracking.
type Skeleton struct {
	skeleton *C.SkeltrackSkeleton
}

// NewSkeleton creates a new Skeleton object. This should always be called
// before use.
func NewSkeleton() Skeleton {

	var skeleton Skeleton

	skeleton.skeleton = C.skeltrack_skeleton_new()

	return skeleton
}

// TrackJoints attempts to find the joints of a human given the depth
// information, width, and height. The depth data should be in millimeters.
// The width and height values are the width and height of the depth buffer.
// Only joints present in the returned map were detected. Keep in mind that
// the default dimension reduction value is 16, so TrackJoints expects the
// depth data to have been shrunk to a 16th of it's original size.
func (skeleton *Skeleton) TrackJoints(depth []uint16, width, height int) (map[JointID]Joint, error) {

	if len(depth) != width*height {
		return make(map[JointID]Joint), errors.New("width and height values do not match depth buffer size")
	}

	var err *C.GError

	sliceHeader := *(*reflect.SliceHeader)(unsafe.Pointer(&depth))

	list := C.skeltrack_skeleton_track_joints_sync(skeleton.skeleton,
		(*C.guint16)(unsafe.Pointer(sliceHeader.Data)), C.guint(width), C.guint(height), nil, &err)

	if err != nil {
		return make(map[JointID]Joint), errors.New(C.GoString((*C.char)(err.message)))
	}

	joints := make(map[JointID]Joint)

	if list == nil {
		return joints, nil
	}

	head := C.skeltrack_joint_list_get_joint(list, C.SKELTRACK_JOINT_ID_HEAD)
	if head != nil {
		joints[JointHead] = Joint{head}
	}
	leftShoulder := C.skeltrack_joint_list_get_joint(list, C.SKELTRACK_JOINT_ID_LEFT_SHOULDER)
	if leftShoulder != nil {
		joints[JointLeftShoulder] = Joint{leftShoulder}
	}
	rightShoulder := C.skeltrack_joint_list_get_joint(list, C.SKELTRACK_JOINT_ID_RIGHT_SHOULDER)
	if rightShoulder != nil {
		joints[JointRightShoulder] = Joint{rightShoulder}
	}
	leftElbow := C.skeltrack_joint_list_get_joint(list, C.SKELTRACK_JOINT_ID_LEFT_ELBOW)
	if leftElbow != nil {
		joints[JointLeftElbow] = Joint{leftElbow}
	}
	rightElbow := C.skeltrack_joint_list_get_joint(list, C.SKELTRACK_JOINT_ID_RIGHT_ELBOW)
	if rightElbow != nil {
		joints[JointRightElbow] = Joint{rightElbow}
	}
	leftHand := C.skeltrack_joint_list_get_joint(list, C.SKELTRACK_JOINT_ID_LEFT_HAND)
	if leftHand != nil {
		joints[JointLeftHand] = Joint{leftHand}
	}
	rightHand := C.skeltrack_joint_list_get_joint(list, C.SKELTRACK_JOINT_ID_RIGHT_HAND)
	if rightHand != nil {
		joints[JointRightHand] = Joint{rightHand}
	}

	return joints, nil
}

// FocusPoint returns the focus point the Skeleton object is currently using
// in millimeters. The focus point is where the tracking starts. The default
// values are x = 0, y = 0, z = 1000.
func (skeleton *Skeleton) FocusPoint() (int, int, int) {

	var x, y, z C.gint

	C.skeltrack_skeleton_get_focus_point(skeleton.skeleton, &x, &y, &z)

	return int(x), int(y), int(z)
}

// SetFocusPoint changes the focus point in millimeters used by the Skeleton
// object. The focus point is where the tracking starts.
func (skeleton *Skeleton) SetFocusPoint(x, y, z int) {

	C.skeltrack_skeleton_set_focus_point(skeleton.skeleton, C.gint(x), C.gint(y), C.gint(z))
}

// DimensionReduction returns the current dimension reduction value. Skeltrack
// uses this to return correct coordinates even when depth data is scaled down.
// The default value is 16.
func (skeleton *Skeleton) DimensionReduction() uint {

	var value C.uint

	C.getObjParamUint((C.gpointer)(skeleton.skeleton), C.CString("dimension-reduction"), &value)

	return uint(value)
}

// SetDimensionReduction sets the current dimension reduction value. Generally
// speaking, depth data should be reduced in size to avoid long wait times
// during analysis. Skeltrack uses this to return correct coordinates even when
// depth data is scaled down. The allowed values are [1, 1024].
func (skeleton *Skeleton) SetDimensionReduction(value uint) {

	C.setObjParamUint((C.gpointer)(skeleton.skeleton), C.CString("dimension-reduction"), C.uint(value))
}

// Smoothing returns true if smoothing is being applied to joints. Smoothing
// makes joint positions less jittery. Smoothing is on by default.
func (skeleton *Skeleton) Smoothing() bool {

	var value C.bool

	C.getObjParamBool((C.gpointer)(skeleton.skeleton), C.CString("enable-smoothing"), &value)

	return bool(value)
}

// SetSmoothing sets whether smoothing should be applied to joints or not.
// Smoothing makes joint positions less jittery.
func (skeleton *Skeleton) SetSmoothing(value bool) {

	C.setObjParamBool((C.gpointer)(skeleton.skeleton), C.CString("enable-smoothing"), C.bool(value))
}

// ExtremaSphereRadius returns the current radius of the sphere around the
// extremas. The default value is 280.
func (skeleton *Skeleton) ExtremaSphereRadius() uint {

	var value C.uint

	C.getObjParamUint((C.gpointer)(skeleton.skeleton), C.CString("extrema-sphere-radius"), &value)

	return uint(value)
}

// SetExtremaSphereRadius sets the current radius of the sphere around the
// extremas. The allowed values are <= 65535.
func (skeleton *Skeleton) SetExtremaSphereRadius(value uint) {

	C.setObjParamUint((C.gpointer)(skeleton.skeleton), C.CString("extrema-sphere-radius"), C.uint(value))
}

// GraphDistanceThreshold returns the current distance threshold between each
// node and it's neighbors. A node in the graph will only be connected to
// another if they aren't further apart then this value. The default value is
// 150.
func (skeleton *Skeleton) GraphDistanceThreshold() uint {

	var value C.uint

	C.getObjParamUint((C.gpointer)(skeleton.skeleton), C.CString("graph-distance-threshold"), &value)

	return uint(value)
}

// SetGraphDistanceThreshold sets the current distance threshold between each
// node and it's neighbors. A node in the graph will only be connected to
// another if they aren't futher apart then this value. The allowed values are
// [1, 65535].
func (skeleton *Skeleton) SetGraphDistanceThreshold(value uint) {

	C.setObjParamUint((C.gpointer)(skeleton.skeleton), C.CString("graph-distance-threshold"), C.uint(value))
}

// GraphMinimumNumberNodes returns the current minimum number of nodes each of
// the graph's components should have. The default value is 5.
func (skeleton *Skeleton) GraphMinimumNumberNodes() uint {

	var value C.uint

	C.getObjParamUint((C.gpointer)(skeleton.skeleton), C.CString("graph-minimum-number-nodes"), &value)

	return uint(value)
}

// SetGraphMinimumNumberNodes sets the current minimum number of nodes each of
// the graph's components should have. The allowed values are [1, 65535].
func (skeleton *Skeleton) SetGraphMinimumNumberNodes(value uint) {

	C.setObjParamUint((C.gpointer)(skeleton.skeleton), C.CString("graph-minimum-number-nodes"), C.uint(value))
}

// HandsMinimumDistance returns the current minimum distance that hands should
// be from the shoulder in millimeters. The default value is 525.
func (skeleton *Skeleton) HandsMinimumDistance() uint {

	var value C.uint

	C.getObjParamUint((C.gpointer)(skeleton.skeleton), C.CString("hands-minimum-distance"), &value)

	return uint(value)
}

// SetHandsMinimumDistance sets the current minimum distance that hands should
// be from the shoulder in millimeters. The allowed values are <= 65535.
func (skeleton *Skeleton) SetHandsMinimumDistance(value uint) {

	C.setObjParamUint((C.gpointer)(skeleton.skeleton), C.CString("hands-minimum-distance"), C.uint(value))
}

// JointsPersistency returns the amount of times a previous value should be
// used for a joint when the joint isn't found. The default value is 3.
func (skeleton *Skeleton) JointsPersistency() uint {

	var value C.uint

	C.getObjParamUint((C.gpointer)(skeleton.skeleton), C.CString("joints-persistency"), &value)

	return uint(value)
}

// SetJointsPersistency sets the amount of times a previous value should be
// used for a joint when the joint isn't found. Allowed values are <= 65535.
func (skeleton *Skeleton) SetJointsPersistency(value uint) {

	C.setObjParamUint((C.gpointer)(skeleton.skeleton), C.CString("joints-persistency"), C.uint(value))
}

// ShouldersArcLength returns the legnth of the arc where the shoulders will
// be searched in millimeters. The default value is 250.
func (skeleton *Skeleton) ShouldersArcLength() uint {

	var value C.uint

	C.getObjParamUint((C.gpointer)(skeleton.skeleton), C.CString("shoulders-arc-length"), &value)

	return uint(value)
}

// SetShouldersArcLength sets the length of the arc where the shoulders will be
// searched in millimeters. Allowed values are [1, 65535].
func (skeleton *Skeleton) SetShouldersArcLength(value uint) {

	C.setObjParamUint((C.gpointer)(skeleton.skeleton), C.CString("shoulders-arc-length"), C.uint(value))
}

// ShouldersArcStartPoint returns the start point of the shoulder searching
// arc. The default value is 120.
func (skeleton *Skeleton) ShouldersArcStartPoint() uint {

	var value C.uint

	C.getObjParamUint((C.gpointer)(skeleton.skeleton), C.CString("shoulders-arc-start-point"), &value)

	return uint(value)
}

// SetShouldersArcStartPoint sets the start point of the shoulder searching
// arc. The allowed values are [1, 65535].
func (skeleton *Skeleton) SetShouldersArcStartPoint(value uint) {

	C.setObjParamUint((C.gpointer)(skeleton.skeleton), C.CString("shoulders-arc-start-point"), C.uint(value))
}

// ShouldersCircumferenceRadius returns the radius of the circumference from
// the head to the shoulders in millimeters. The default value is 290.
func (skeleton *Skeleton) ShouldersCircumferenceRadius() uint {

	var value C.uint

	C.getObjParamUint((C.gpointer)(skeleton.skeleton), C.CString("shoulders-circumference-radius"), &value)

	return uint(value)
}

// SetShouldersCircumferenceRadius sets the radius of the circumference from
// the head to the shoulders in millimeters. Allowed values are [1, 65535].
func (skeleton *Skeleton) SetShouldersCircumferenceRadius(value uint) {

	C.setObjParamUint((C.gpointer)(skeleton.skeleton), C.CString("shoulders-circumference-radius"), C.uint(value))
}

// ShouldersSearchStep returns the step used when searching for shoulders. The
// default value is 0.01.
func (skeleton *Skeleton) ShouldersSearchStep() float32 {

	var value C.float

	C.getObjParamFloat((C.gpointer)(skeleton.skeleton), C.CString("shoulders-search-step"), &value)

	return float32(value)
}

// SetShouldersSearchStep sets the step used when searching for shoulders. The
// allowed values are [0.01, 3.14159].
func (skeleton *Skeleton) SetShouldersSearchStep(value float32) {

	C.setObjParamFloat((C.gpointer)(skeleton.skeleton), C.CString("shoulders-search-step"), C.float(value))
}

// SmoothingFactor returns the current smoothing factor being used if smoothing
// is turned on. A higher value will have smoother results. The default value
// is 0.5.
func (skeleton *Skeleton) SmoothingFactor() float32 {

	var value C.float

	C.getObjParamFloat((C.gpointer)(skeleton.skeleton), C.CString("smoothing-factor"), &value)

	return float32(value)
}

// SetSmoothingFactor sets the current smoothing factor being used if smoothing
// is turned on. A higher value will have smoothe results. Allowed values are
// [0, 1].
func (skeleton *Skeleton) SetSmoothingFactor(value float32) {

	C.setObjParamFloat((C.gpointer)(skeleton.skeleton), C.CString("smoothing-factor"), C.float(value))
}

// TorsoMinimumNumberNodes returns the current minimun number of nodes required
// for a component to be considered a torso. The default value is 16.
func (skeleton *Skeleton) TorsoMinimumNumberNodes() float32 {

	var value C.float

	C.getObjParamFloat((C.gpointer)(skeleton.skeleton), C.CString("torso-minimum-number-nodes"), &value)

	return float32(value)
}

// SetTorsoMinimumNumberNodes sets the current minimum number of nodes required
// for a component to be considered a torso. Allowed values are [0, 65535].
func (skeleton *Skeleton) SetTorsoMinimumNumberNodes(value float32) {

	C.setObjParamFloat((C.gpointer)(skeleton.skeleton), C.CString("torso-minimum-number-nodes"), C.float(value))
}
