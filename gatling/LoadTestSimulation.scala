package computerdatabase

import io.gatling.core.Predef._
import io.gatling.http.Predef._
import scala.concurrent.duration._

import scala.util.parsing.json._
import scala.util.Random
import java.security.SecureRandom

class LoadTestSimulation extends Simulation {
  def displayParam(userID: String, name: String, secretKey: String): String = {
    val orderID = new Random(new SecureRandom()).alphanumeric.take(10).mkString
    val json =
      s"""{
        "user_id": "${userID}",
        "name": "${name}",
        "secret_key": "${secretKey}"
      }"""
    return json
  }
  val serverURI = "http://localhost:8080/v1"
  val serverHttpProtocol = http.baseUrl(serverURI)
  val sc = scenario("Scenario")
    .exec(
      http("scenario_1_create_user")
        .post(serverURI + "/users")
        .check(
          status.is(201),
          jsonPath("$.id").find.saveAs("user_id")
        )
    ).exitHereIfFailed
    .exec(
      http("scenario_2_refer_user")
        .get(session => serverURI+"/users/"+session("user_id").as[String])
        .check(
          status.is(200),
          jsonPath("$.id").find.saveAs("user_id"),
          jsonPath("$.name").find.saveAs("name"),
          jsonPath("$.secret_key").find.saveAs("secret_key")
        ),
    ).exitHereIfFailed
    .exec(
      http("scenario_3_redirect_status")
        .post(serverURI+"/display")
        .body(StringBody(session => displayParam(session("user_id").as[String], session("name").as[String], session("secret_key").as[String]))).asJson
        .disableFollowRedirect
        .check(status.is(301))
    ).exitHereIfFailed
    .exec(
      http("scenario_3_redirect_page")
        .post(serverURI+"/display")
        .body(StringBody(session => displayParam(session("user_id").as[String], session("name").as[String], session("secret_key").as[String]))).asJson
        .check(status.is(200))
    ).exitHereIfFailed

  setUp(
    sc.inject(
      atOnceUsers(10),
      rampUsers(20) during(5 seconds),
    ).protocols(serverHttpProtocol),
  )
}