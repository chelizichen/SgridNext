package org.chelizichen.sgridjava.demo.sgrid;
import org.springframework.stereotype.Component;


@Component
public class Constant {

    public final static String SGRID_LOG_DIR = "SGRID_LOG_DIR";
    public final static String SGRID_CONF_DIR = "SGRID_CONF_DIR";
    public final static String SGRID_PACKAGE_DIR = "SGRID_PACKAGE_DIR";
    public final static String SGRID_SERVANT_DIR = "SGRID_SERVANT_DIR";
    public final static String SGRID_TARGET_PORT = "SGRID_TARGET_PORT";
    public final static String SGRID_TARGET_HOST = "SGRID_TARGET_HOST";

    public static boolean IsProduction() {
        return System.getenv(SGRID_TARGET_PORT) != null;
    }

    public static String getConfDir() {
        return System.getenv(SGRID_CONF_DIR);
    }

}
