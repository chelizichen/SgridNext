package org.chelizichen.sgridjava.demo.sgrid.init;

import org.springframework.context.annotation.Bean;
import org.springframework.stereotype.Component;

@Component
public class SgridStarter {
    @Bean
    public ServletContainer servletContainer() {
        System.out.println("[Sgrid-Java] [info] init servletContainer ");
        return new ServletContainer();
    }
}
