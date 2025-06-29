class CloudManager {
  constructor() {
    this.cloudImages = [
      "cloud1.png",
      "cloud2.png",
      "cloud3.png",
      "cloud4.png",
      "cloud5.png",
      "cloud6.png",
      "cloud7.png",
      "cloud8.png",
      "cloud9.png",
    ];
    this.cloudSizes = ["small", "medium", "large"];
    this.cloudsContainer = null;
    this.cloudCount = 12; // Total number of clouds
  }

  init() {
    this.createCloudsContainer();
    this.generateClouds();
    console.log("🌤️ Cloud system initialized with", this.cloudCount, "clouds");
  }

  createCloudsContainer() {
    const existing = document.querySelector(".clouds-container");
    if (existing) {
      existing.remove();
    }

    this.cloudsContainer = document.createElement("div");
    this.cloudsContainer.className = "clouds-container";
    document.body.insertBefore(this.cloudsContainer, document.body.firstChild);
  }

  generateClouds() {
    for (let i = 1; i <= this.cloudCount; i++) {
      const cloud = this.createCloudElement(i);
      this.cloudsContainer.appendChild(cloud);
    }
  }

  createCloudElement(index) {
    const cloud = document.createElement("img");

    const randomImage =
      this.cloudImages[Math.floor(Math.random() * this.cloudImages.length)];
    cloud.src = `/static/sprites/${randomImage}`;
    cloud.alt = `Cloud ${index}`;

    const randomSize =
      this.cloudSizes[Math.floor(Math.random() * this.cloudSizes.length)];

    cloud.className = `cloud cloud-${randomSize} cloud-${index}`;

    const randomTop = Math.random() * 80 + 5;
    cloud.style.top = `${randomTop}%`;

    const randomOpacity = 0.4 + Math.random() * 0.4;
    cloud.style.opacity = randomOpacity;

    const randomDuration = 35 + Math.random() * 40;
    cloud.style.animationDuration = `${randomDuration}s`;

    const randomDelay = -Math.random() * randomDuration;
    cloud.style.animationDelay = `${randomDelay}s`;

    cloud.onerror = () => {
      console.warn(`Cloud image not found: ${randomImage}`);
      cloud.style.display = "none";
    };

    return cloud;
  }

  addCloud() {
    const newIndex = this.cloudCount + 1;
    const cloud = this.createCloudElement(newIndex);
    this.cloudsContainer.appendChild(cloud);
    this.cloudCount++;
  }

  removeClouds() {
    if (this.cloudsContainer) {
      this.cloudsContainer.remove();
    }
  }

  toggleAnimation() {
    const clouds = document.querySelectorAll(".cloud");
    clouds.forEach((cloud) => {
      if (cloud.style.animationPlayState === "paused") {
        cloud.style.animationPlayState = "running";
      } else {
        cloud.style.animationPlayState = "paused";
      }
    });
  }
}

document.addEventListener("DOMContentLoaded", function () {
  setTimeout(() => {
    window.cloudManager = new CloudManager();
    window.cloudManager.init();
  }, 100);
});

// Make CloudManager available globally for debugging
window.CloudManager = CloudManager;
