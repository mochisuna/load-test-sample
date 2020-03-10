package computerdatabase

import io.gatling.core.Predef._
import io.gatling.http.Predef._
import scala.concurrent.duration._

import scala.util.parsing.json._
import scala.util.Random
import java.security.SecureRandom

class LoadTestSimulation extends Simulation {
  val serverURI = "http://localhost:28080/v1"
  val serverHttpProtocol = http.baseUrl(serverURI)

  object User {
    def create() = exec(
      http("create_user")
        .post(f"$serverURI/users")
        .check(
          status.is(201),
          jsonPath("$.id").find.saveAs("user_id"),
          jsonPath("$.name").find.saveAs("user_name")
        )
    ).exitHereIfFailed

    def refer(userID: String, userName: String) = exec(
      http("refer_user")
        .get(f"$serverURI/users/$userID")
        .check(
          status.is(200),
          jsonPath("$.id").find.is(userID),
          jsonPath("$.name").find.is(userName),
          jsonPath("$.secret_key").find.saveAs("secret_key")
        ),
    ).exitHereIfFailed
  }
  object Page{
    private val param = (userID: String, name: String, secretKey: String) => {
      StringBody(
        s"""{
          "user_id": "$userID",
          "name": "$name",
          "secret_key": "$secretKey"
        }"""
      )
    }
    def redirectStatus(userID: String, name: String, secretKey: String) = exec(
      http("redirect_status")
        .post(f"$serverURI/display")
        .body(param(userID, name, secretKey)).asJson
        .disableFollowRedirect
        .check(status.is(301))
    ).exitHereIfFailed

    def display(userID: String, name: String, secretKey: String) = exec(
      http("redirect_page")
        .post(f"$serverURI/display")
        .body(param(userID, name, secretKey)).asJson
        .check(status.is(200))
    ).exitHereIfFailed
  }
  val sc = scenario("Scenario")
    .exec(User.create())
    .exec(User.refer("${user_id}", "${user_name}"))
    .exec(Page.redirectStatus("${user_id}", "${user_name}", "${secret_key}"))
    .exec(Page.display("${user_id}", "${user_name}", "${secret_key}"))

  setUp(
    sc.inject(
      atOnceUsers(10),
      rampUsers(20) during(5 seconds),
    ).protocols(serverHttpProtocol),
  )
}