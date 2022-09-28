package spotandmargin

import "github.com/chuckpreslar/emission"

//On adds a listener to a specific event
func (b *GateWS) On(event interface{}, listener interface{}) *emission.Emitter {
	return b.emitter.On(event, listener)
}

//Emit emits an event
func (b *GateWS) Emit(event interface{}, arguments ...interface{}) *emission.Emitter {
	return b.emitter.Emit(event, arguments...)
}

//Off removes a listener for an event
func (b *GateWS) Off(event interface{}, listener interface{}) *emission.Emitter {
	return b.emitter.Off(event, listener)
}
