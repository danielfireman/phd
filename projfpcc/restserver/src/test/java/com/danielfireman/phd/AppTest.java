package com.danielfireman.phd;

import org.junit.Test;

public class AppTest extends BaseTest {

  @Test
  public void index() throws Exception {
    server.post("/echo")
    	.body("\"content\":\"Hi\"", "application/json;charset=UTF-8")
    	.expect(200);
  }
}
