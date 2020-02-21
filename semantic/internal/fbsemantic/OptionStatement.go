// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package fbsemantic

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type OptionStatement struct {
	_tab flatbuffers.Table
}

func GetRootAsOptionStatement(buf []byte, offset flatbuffers.UOffsetT) *OptionStatement {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &OptionStatement{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *OptionStatement) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *OptionStatement) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *OptionStatement) Loc(obj *SourceLocation) *SourceLocation {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(SourceLocation)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *OptionStatement) AssignmentType() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *OptionStatement) MutateAssignmentType(n byte) bool {
	return rcv._tab.MutateByteSlot(6, n)
}

func (rcv *OptionStatement) Assignment(obj *flatbuffers.Table) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		rcv._tab.Union(obj, o)
		return true
	}
	return false
}

func OptionStatementStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func OptionStatementAddLoc(builder *flatbuffers.Builder, loc flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(loc), 0)
}
func OptionStatementAddAssignmentType(builder *flatbuffers.Builder, assignmentType byte) {
	builder.PrependByteSlot(1, assignmentType, 0)
}
func OptionStatementAddAssignment(builder *flatbuffers.Builder, assignment flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(assignment), 0)
}
func OptionStatementEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}