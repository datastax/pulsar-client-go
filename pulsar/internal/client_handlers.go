// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package internal

import "sync"

// ClientHandlerMap is a simple concurrent-safe map for the client type
type ClientHandlers struct {
	handlers map[Closable]bool
	sync.RWMutex
}

func NewClientHandlers() ClientHandlers {
	return ClientHandlers{
		handlers: map[Closable]bool{},
	}
}
func (h *ClientHandlers) Add(c Closable) {
	h.Lock()
	defer h.Unlock()
	h.handlers[c] = true
}

func (h *ClientHandlers) Del(c Closable) {
	h.Lock()
	defer h.Unlock()
	delete(h.handlers, c)
}

func (h *ClientHandlers) Val(c Closable) bool {
	h.RLock()
	defer h.RUnlock()
	return h.handlers[c]
}

func (h *ClientHandlers) Close() {
	h.Lock()
	handlers := make([]Closable, 0, len(h.handlers))
	for handler := range h.handlers {
		handlers = append(handlers, handler)
	}
	h.Unlock()

	for _, handler := range handlers {
		handler.Close()
	}
}
