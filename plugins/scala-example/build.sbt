name := "scala-example-plugin"

libraryDependencies ++= Seq(
  "io.grpc" % "grpc-netty" % scalapb.compiler.Version.grpcJavaVersion,
  "com.thesamet.scalapb" %% "scalapb-runtime-grpc" % scalapb.compiler.Version.scalapbVersion
)

scalaVersion := "2.12.8"

PB.targets in Compile := Seq(
  scalapb.gen() -> (sourceManaged in Compile).value
)
PB.protoSources in Compile := Seq(baseDirectory.value / ".." / ".." / "proto")

cancelable in Global := true

enablePlugins(JavaAppPackaging)
enablePlugins(WindowsPlugin)
