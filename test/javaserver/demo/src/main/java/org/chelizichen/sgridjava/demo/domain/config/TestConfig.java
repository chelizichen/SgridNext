package org.chelizichen.sgridjava.demo.domain.config;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Configuration;

@Configuration
public class TestConfig {
    @Value("${app.serverName}")
    private String serverName;

    public String getServerName() {
        return serverName;
    }
}
