package main

import (
	"context"
	"fmt"
	"go-plugin-example/proto"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/hashicorp/go-plugin"
)

func main() {
	log.SetOutput(ioutil.Discard)

	// ハンドシェイクの内容　クライアント/サーバのプロトコルバージョンの差異が無いかをチェックする
	var Handshake = plugin.HandshakeConfig{
		// This isn't required when using VersionedPlugins
		ProtocolVersion:  1,
		MagicCookieKey:   "BASIC_PLUGIN",
		MagicCookieValue: "hello",
	}

	// プラグイン名と実装の対応付　プラグインは、サーバかクライアント開始する
	var PluginMap = map[string]plugin.Plugin{
		"kv_grpc": &KVGRPCPlugin{},
	}

	// 環境変数から受け取ったプロセス名でExecし、標準出力でのハンドシェイクを行う
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  Handshake,
		Plugins:          PluginMap,
		Cmd:              exec.Command(os.Getenv("PLUGIN")),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})
	defer client.Kill()

	// クライアントと接続する
	rpcClient, err := client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// 名前を指定して、プラグインのクライアントインスタンスを取り出す
	raw, err := rpcClient.Dispense("kv_grpc")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// 任意の型にキャストしてクライアントを操作する
	service := raw.(SayHello)
	result, err := service.Hello("hoge")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	fmt.Printf("message: %s\n", result)
	os.Exit(0)
}

// プラグインの実装　サーバorクライアントを開始する
type KVGRPCPlugin struct {
	plugin.Plugin
}

func (p *KVGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	return nil
}

func (p *KVGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: com_example_protos.NewGreeterClient(c)}, nil
}

// クライアントのインターフェース
type SayHello interface {
	Hello(string) (string, error)
}

// クライアントの実装　gRPCで実際に通信を行う
type GRPCClient struct {
	client com_example_protos.GreeterClient
}

func (c *GRPCClient) Hello(name string) (string, error) {
	req := &com_example_protos.HelloRequest{
		Name: name,
	}
	res, err := c.client.SayHello(context.TODO(), req)
	return res.Message, err
}
