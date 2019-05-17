import com.example.protos.health.{HealthCheckRequest, HealthCheckResponse, HealthGrpc}
import com.example.protos.hello.{GreeterGrpc, HelloReply, HelloRequest}

import scala.concurrent.Future

// サービスのヘルスチェックを行うための実装　go-plugin側で定義されているサービスを実装している
class HealthCheckImpl extends HealthGrpc.Health {
  override def check(request: HealthCheckRequest): Future[HealthCheckResponse] = {
    val res = HealthCheckResponse(HealthCheckResponse.ServingStatus.SERVING)
    Future.successful(res)
  }
}

// 自分で定義したサービス
class GreeterImpl extends GreeterGrpc.Greeter {
  override def sayHello(req: HelloRequest) = {
    val reply = HelloReply(message = "Hello " + req.name)
    Future.successful(reply)
  }
}
