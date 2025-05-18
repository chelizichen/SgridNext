package org.chelizichen.sgridjava.demo.sgrid;


import java.nio.file.Path;
import java.nio.file.Paths;

public class Constant {

    public final static String SGRID_LOG_DIR = "SGRID_LOG_DIR";
    public final static String SGRID_CONF_DIR = "SGRID_CONF_DIR";
    public final static String SGRID_PACKAGE_DIR = "SGRID_PACKAGE_DIR";
    public final static String SGRID_SERVANT_DIR = "SGRID_SERVANT_DIR";
    public final static String SGRID_TARGET_PORT = "SGRID_TARGET_PORT";
    public final static String SGRID_TARGET_HOST = "SGRID_TARGET_HOST";

    public static boolean IsProduction() {
        System.out.println("System.getenv(SGRID_TARGET_PORT) >> " + System.getenv(SGRID_TARGET_PORT));
        return System.getenv(SGRID_TARGET_PORT) != null;
    }

    public static String getConfDir() {
        return System.getenv(SGRID_CONF_DIR);
    }

    public static String getLocalResourcesDir() {
        if (Constant.IsProduction()) {
            String currentDir = System.getProperty("user.dir");
            System.out.println("Current working directory: " + currentDir);
            Path path = Paths.get(currentDir, "target/classes/conf/");
            return path.toString();
        }
        return "src/main/resources/conf";
    }

}
