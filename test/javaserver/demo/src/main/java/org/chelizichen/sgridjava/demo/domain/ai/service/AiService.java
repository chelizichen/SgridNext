package org.chelizichen.sgridjava.demo.domain.ai.service;

import org.springframework.stereotype.Service;

@Service
public class AiService {
    public Integer fibonacci(int n) {
        if (n <= 1) {
            return n;
        }
        return fibonacci(n - 1) + fibonacci(n - 2);
    }
}
