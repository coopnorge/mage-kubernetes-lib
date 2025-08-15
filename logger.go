package magekubernetes

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)

func infof(format string, v ...any)  { logger.Printf("[INFO] "+format, v...) }
func warnf(format string, v ...any)  { logger.Printf("[WARN] "+format, v...) }
func errorf(format string, v ...any) { logger.Printf("[ERROR] "+format, v...) }
func debugf(format string, v ...any) { logger.Printf("[DEBUG] "+format, v...) }
