#!/usr/bin/env python3
import time
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import Select
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from webdriver_manager.chrome import ChromeDriverManager
from selenium.webdriver.chrome.service import Service

# Initialize driver
service = Service(ChromeDriverManager().install())
driver = webdriver.Chrome(service=service)
driver.maximize_window()

BASE_URL = "https://gleaming-pudding-326d57.netlify.app"
PASSWORD = "meet21"

def log(msg):
    print(f"[TEST] {msg}")

def test_login():
    log("Starting login test...")
    driver.get(f"{BASE_URL}/")
    time.sleep(2)
    
    log(f"Current URL: {driver.current_url}")
    log(f"Page title: {driver.title}")
    
try:
    test_login()
    time.sleep(5)
    
    log("Taking screenshot...")
    driver.save_screenshot("/tmp/ui_test_1.png")
    
    log("Screenshot saved to /tmp/ui_test_1.png")
    
finally:
    driver.quit()
