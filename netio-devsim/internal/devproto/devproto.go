/*
Copyright © 2021 Ci4Rail GmbH
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package devproto

import (
	"github.com/ci4rail/firmware-ota/pkg/socket"
	"github.com/golang/protobuf/proto"
)

type DevProto struct {
	Listener *socket.Listener
	Socket   *socket.Socket
}

func NewDevProto(port string) (*DevProto, error) {
	l, err := socket.NewListener(port)
	if err != nil {
		return &DevProto{}, err
	}
	return &DevProto{
		Listener: l,
	}, nil
}

func (p *DevProto) WaitForConnection() error {
	socket, err := socket.WaitForConnect(p.Listener)
	if err != nil {
		return err
	}
	p.Socket = socket
	return nil
}

func (p *DevProto) WriteMessage(m proto.Message) error {
	payload, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return p.Socket.Write(payload)
}

func (p *DevProto) ReadMessage(m proto.Message) error {
	payload, err := p.Socket.Read()
	if err != nil {
		return err
	}

	return proto.Unmarshal(payload, m)
}

func (p *DevProto) Close() {
	p.Socket.Close()
}
