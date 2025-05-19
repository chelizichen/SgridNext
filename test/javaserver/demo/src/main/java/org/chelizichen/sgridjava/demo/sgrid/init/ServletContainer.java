package org.chelizichen.sgridjava.demo.sgrid.init;

import org.chelizichen.sgridjava.demo.sgrid.Constant;
import org.springframework.boot.web.server.WebServerFactoryCustomizer;
import org.springframework.boot.web.servlet.server.ConfigurableServletWebServerFactory;

import java.net.InetAddress;
import java.net.UnknownHostException;

public class ServletContainer implements WebServerFactoryCustomizer<ConfigurableServletWebServerFactory> {

    @Override
    public void customize(ConfigurableServletWebServerFactory factory) {
        try {
            System.out.println("[Sgrid-Java] [info] start Init Sgrid servletContainer ");
            System.out.println("[Sgrid-Java] [info] get sgrid remote configuration ");
            if (Constant.IsProduction()) {
                System.out.println("[Sgrid-Java] [info] get sgrid remote configuration ");
                String port = System.getenv(Constant.SGRID_TARGET_PORT);
                String host = System.getenv(Constant.SGRID_TARGET_HOST);
                factory.setPort(Integer.parseInt(port));
                factory.setAddress(InetAddress.getByName(host));
            }
            System.out.println("[Sgrid-Java] [info] End Init Sgrid servletContainer ");

        } catch (UnknownHostException e) {
            throw new RuntimeException(e);
        }
    }
}
