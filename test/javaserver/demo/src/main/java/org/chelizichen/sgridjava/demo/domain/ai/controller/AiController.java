package org.chelizichen.sgridjava.demo.domain.ai.controller;

import org.chelizichen.sgridjava.demo.domain.ai.service.AiService;
import org.chelizichen.sgridjava.demo.domain.config.TestConfig;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/ai")
public class AiController {
    @Autowired
    private TestConfig testConfig;
    @Autowired
    private AiService aiService;

    @RequestMapping(value = "/hello", method = RequestMethod.POST)
    public String hello() {
        return "Hello World";
    }

    @RequestMapping(value = "/fib", method = RequestMethod.GET)
    public Integer fib(@RequestParam("n") Integer n) {
        System.out.println(testConfig.getServerName());
        // 计算斐波那契数列
        Integer fibonacci = aiService.fibonacci(n);
        return fibonacci;
    }
}
