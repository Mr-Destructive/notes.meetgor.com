#!/usr/bin/env python3
import time
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import Select, WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from webdriver_manager.chrome import ChromeDriverManager
from selenium.webdriver.chrome.service import Service

service = Service(ChromeDriverManager().install())
driver = webdriver.Chrome(service=service)
driver.maximize_window()

BASE_URL = "https://gleaming-pudding-326d57.netlify.app"
PASSWORD = "meet21"

def log(msg):
    print(f"[TEST] {msg}")

def screenshot(name):
    driver.save_screenshot(f"/tmp/{name}.png")
    log(f"Screenshot: {name}")

try:
    # Test 1: Check if dashboard loads
    log("Test 1: Loading dashboard...")
    driver.get(f"{BASE_URL}/")
    time.sleep(3)
    screenshot("01_initial_load")
    
    log(f"Title: {driver.title}")
    log(f"URL: {driver.current_url}")
    
    # Check for admin dashboard or login
    try:
        dashboard = driver.find_element(By.ID, "page-title")
        log(f"✓ Found page title: {dashboard.text}")
    except:
        log("✗ No page-title found, checking for login")
    
    # Test 2: Check sidebar navigation
    try:
        nav_links = driver.find_elements(By.CLASS_NAME, "nav-link")
        log(f"✓ Found {len(nav_links)} navigation links")
        for link in nav_links[:3]:
            log(f"  - {link.text}")
    except Exception as e:
        log(f"✗ No nav-links found: {e}")
    
    # Test 3: Navigate to posts and create new post
    log("Test 2: Navigating to posts...")
    try:
        posts_link = driver.find_element(By.XPATH, "//*[contains(text(), 'Posts')]")
        driver.execute_script("arguments[0].scrollIntoView();", posts_link)
        posts_link.click()
        time.sleep(2)
        screenshot("02_posts_page")
    except Exception as e:
        log(f"✗ Could not click Posts: {e}")
    
    # Test 4: Click new post button
    log("Test 3: Creating new post...")
    try:
        new_btn = driver.find_element(By.XPATH, "//*[contains(text(), 'New Post')]")
        new_btn.click()
        time.sleep(2)
        screenshot("03_new_post_form")
    except Exception as e:
        log(f"✗ Could not find New Post button: {e}")
    
    # Test 5: Fill in post form
    log("Test 4: Filling in post form...")
    try:
        # Select post type
        type_select = Select(driver.find_element(By.ID, "post-type"))
        type_select.select_by_value("article")
        log("✓ Selected 'article' type")
        time.sleep(1)
        screenshot("04_article_type_selected")
        
        # Fill title
        title_input = driver.find_element(By.ID, "post-title")
        title_input.send_keys("Test Article")
        log("✓ Entered title")
        
        # Fill slug
        slug_input = driver.find_element(By.ID, "post-slug")
        slug_input.send_keys("test-article")
        log("✓ Entered slug")
        
        # Fill content
        content_input = driver.find_element(By.ID, "post-content")
        content_input.send_keys("This is test content for the article.")
        log("✓ Entered content")
        
        screenshot("05_form_filled")
        
    except Exception as e:
        log(f"✗ Error filling form: {e}")
    
    # Test 6: Save as draft
    log("Test 5: Saving as draft...")
    try:
        save_draft_btn = driver.find_element(By.XPATH, "//*[text()='💾 Save as Draft']")
        save_draft_btn.click()
        log("✓ Clicked Save as Draft")
        time.sleep(3)
        screenshot("06_after_save_draft")
        
        # Check for success message
        try:
            message = driver.find_element(By.ID, "editor-message")
            log(f"Message: {message.text}")
        except:
            log("No editor message visible")
            
    except Exception as e:
        log(f"✗ Error saving draft: {e}")
    
    # Test 7: Check if post appears in list
    log("Test 6: Checking posts list...")
    time.sleep(2)
    try:
        posts_list = driver.find_elements(By.CSS_SELECTOR, "table tbody tr")
        log(f"✓ Found {len(posts_list)} posts in list")
        for post in posts_list[:3]:
            log(f"  - {post.text[:80]}")
    except Exception as e:
        log(f"✗ Could not view posts list: {e}")
    
    screenshot("07_posts_list")

finally:
    driver.quit()
    log("Test complete")
