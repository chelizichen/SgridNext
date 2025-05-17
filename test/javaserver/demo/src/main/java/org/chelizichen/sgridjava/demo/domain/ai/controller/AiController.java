package org.chelizichen.sgridjava.demo.domain.ai.controller;

import org.chelizichen.sgridjava.demo.domain.config.TestConfig;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/ai")
public class AiController {
    @Autowired
    private TestConfig testConfig;

    @RequestMapping(value = "/hello", method = RequestMethod.POST)
    public String hello() {
        System.out.println(testConfig.getServerName());
        return "Hello World";
    }
}
