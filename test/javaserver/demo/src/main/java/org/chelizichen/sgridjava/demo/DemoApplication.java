package org.chelizichen.sgridjava.demo;


import org.chelizichen.sgridjava.demo.sgrid.EnableSgridServer;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
@EnableSgridServer
public class DemoApplication {


    public static void main(String[] args) {
        SpringApplication.run(DemoApplication.class, args);
    }

}
