package com.danielfireman.phd;

import org.junit.Test;

public class AppTest extends BaseTest {

  @Test
  public void index() throws Exception {
    server.get("/msg").expect(200);
  }
}
