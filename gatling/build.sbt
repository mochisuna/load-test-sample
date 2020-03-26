enablePlugins(GatlingPlugin)

scalaVersion := "2.12.10"

scalacOptions := Seq(
  "-encoding", "UTF-8", "-target:jvm-1.8", "-deprecation",
  "-feature", "-unchecked", "-language:implicitConversions", "-language:postfixOps")

libraryDependencies += "io.gatling.highcharts" % "gatling-charts-highcharts" % "3.3.1" % "test"
libraryDependencies += "io.gatling"            % "gatling-test-framework"    % "3.3.1" % "test"
libraryDependencies += "org.scalaj"            %% "scalaj-http"              % "2.4.2"
