import com.example.protos.health.HealthGrpc
import com.example.protos.hello.GreeterGrpc
import io.grpc.ServerBuilder

import scala.concurrent.ExecutionContext

object Main {
  def main(args: Array[String]): Unit = {
    val executionContext = ExecutionContext.global
    val port = 50001

    // HealthCheck/自分のサービスを登録したgRPCサーバを開始する
    val builder = ServerBuilder
      .forPort(port)
    builder
      .addService(HealthGrpc.bindService(new HealthCheckImpl, executionContext))
      .addService(GreeterGrpc.bindService(new GreeterImpl, executionContext))
    val server = builder
      .build
      .start

    // Ctrl-cで終了するようにする
    sys.addShutdownHook {
      server.shutdown()
    }

    // バージョン・アドレスを標準出力に出力し、ネゴシエーションを行う
    println(s"1|1|tcp|127.0.0.1:$port|grpc")

    // サーバが終了する前mainスレッドをブロックする
    server.awaitTermination()
  }
}
