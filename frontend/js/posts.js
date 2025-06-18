let currentCategory = null;

// Fetch and display categories
async function fetchCategories() {
  try {
    const response = await fetch("/categories");
    const categories = await response.json();
    displayCategories(categories);
  } catch (error) {
    console.error("Error fetching categories:", error);
  }
}

// Display categories in the categories grid
function displayCategories(categories) {
  const categoriesList = document.getElementById("categories-list");
  categoriesList.innerHTML = "";

  // Add "All Categories" option
  const allCategoriesElement = document.createElement("div");
  allCategoriesElement.className = "category-card active";
  allCategoriesElement.innerHTML = `
        <h3>All  Posts</h3>
    `;
  allCategoriesElement.onclick = () => {
    // Remove active class from all category cards
    document
      .querySelectorAll(".category-card")
      .forEach((card) => card.classList.remove("active"));
    // Add active class to this card
    allCategoriesElement.classList.add("active");
    currentCategory = null;
    fetchPosts();
  };
  categoriesList.appendChild(allCategoriesElement);
  
  categories.forEach((category) => {
    const categoryElement = document.createElement("div");
    categoryElement.className = "category-card";
    categoryElement.innerHTML = `
            <h3>${category.name}</h3>
        `;
    categoryElement.onclick = () => {
      // Remove active class from all category cards
      document
        .querySelectorAll(".category-card")
        .forEach((card) => card.classList.remove("active"));
      // Add active class to this card
      categoryElement.classList.add("active");
      filterPostsByCategory(category.name);
    };
    categoriesList.appendChild(categoryElement);
  });
}

// Fetch and display posts
async function fetchPosts() {
  try {
    const url = currentCategory
      ? `/posts?category=${encodeURIComponent(currentCategory)}`
      : "/posts";
    const response = await fetch(url);
    const data = await response.json();
    currentUserUUID = data.current_user_uuid;
    displayPosts(data.posts);
    const filterInfo = document.getElementById("post-filter-info");
    if (filterInfo) {
      if (currentCategory === "ReactedPosts") {
        filterInfo.textContent =
          "Showing posts you've reacted to (liked or disliked)";
      } else if (currentCategory) {
        filterInfo.textContent = `Showing posts in category: ${currentCategory}`;
      } else {
        filterInfo.textContent = "Showing all posts";
      }
    }
  } catch (error) {
    console.error("Error fetching posts:", error);
  }
}

// Display posts in the posts grid
function displayPosts(posts) {
  const postsList = document.getElementById("posts-list");
  postsList.innerHTML = "";

  if (!Array.isArray(posts)) {
    postsList.innerHTML = "";
    return;
  }

  posts.forEach((post) => {
    const postElement = document.createElement("div");
    postElement.className = "post-card";
    postElement.setAttribute("data-post-id", post.post_id);
    postElement.innerHTML = `
            <h3>${post.title}</h3>
            <p class="post-content">${post.content}</p>
            ${
              // Check if post has an image_url and display it
              post.image_url
                ? `<div class="post-image-container"><img src="${post.image_url}" alt="Post Image" class="post-image"></div>`
                : ""
            }
            <div class="post-meta">
                <span class="post-author">By ${post.nickname}</span> 
                <span class="post-date">${new Date(
                  post.created_at
                ).toLocaleString()}</span>
            </div>
            <div class="post-categories">
                ${(post.categories || [])
                  .map((cat) => `<span class="category-tag">${cat}</span>`)
                  .join("")}
            </div>
            <div class="post-actions">
                <button onclick="likePost(${
                  post.post_id
                })" class="like-btn">👍 ${post.likes || 0}</button>
            <button onclick="dislikePost(${post.post_id})" class="dislike-btn">
           👎 ${post.dislikes || 0}
             </button>
                <button onclick="toggleComments(${
                  post.post_id
                })" class="comment-btn">💬 Comments (${
             post.comments_count || 0
          })</button>
                ${
                  post.user_uuid === currentUserUUID
                    ? `<button onclick="deletePost(${post.post_id})" class="delete-btn">🗑️ Delete</button>`
                    : ""
                }
            </div>
            <div id="comments-section-${
              post.post_id
            }" class="comments-section" style="display: none;">
                <div id="comments-${
                  post.post_id
                }" class="comments-container"></div>
                <div class="comment-form">
                    <textarea id="comment-input-${
                      post.post_id
                    }" placeholder="Write a comment..." class="comment-input"></textarea>
                    <button type="button" onclick="submitComment(${
                      post.post_id
                    })" class="post-btn">Post Comment</button>
                </div>
            </div>
        `;
    postsList.appendChild(postElement);
  });
}

// Filter posts by category
async function filterPostsByCategory(categoryName) {
  try {
    currentCategory = categoryName;
    const response = await fetch(
      `/posts?category=${encodeURIComponent(categoryName)}`
    );
    const data = await response.json();
    currentUserUUID = data.current_user_uuid;
    displayPosts(data.posts);
  } catch (error) {
    console.error("Error filtering posts:", error);
  }
}

// Populate categories for post creation
async function populateCategoriesForPost() {
  try {
    const response = await fetch("/categories");
    const categories = await response.json();
    const categoriesSelection = document.getElementById("categories-selection");

    categoriesSelection.innerHTML = categories
      .map(
        (category) => `
            <div class="category-checkbox">
                <input type="checkbox" id="category-${category.id}" name="categories" value="${category.id}">
                <label for="category-${category.id}">${category.name}</label>
            </div>
        `
      )
      .join("");
  } catch (error) {
    console.error("Error fetching categories for post:", error);
  }
}

// Event listener for image input change (for preview)
document.addEventListener("DOMContentLoaded", () => {
  const postImageInput = document.getElementById("postImageInput");
  const imagePreview = document.getElementById("imagePreview");
  fetchCategories();
  fetchPosts();

  if (postImageInput && imagePreview) {
    postImageInput.addEventListener("change", function (event) {
      const file = event.target.files[0];
      if (file) {
        // Basic validation for image type and size
        if (!file.type.startsWith("image/")) {
          alert("Please select an image file (JPG, JPEG, PNG, GIF).");
          postImageInput.value = ""; // Clear the input
          imagePreview.style.display = "none";
          return;
        }
        if (file.size > 10 * 1024 * 1024) {
          // 10 MB
          alert("Image file size must not exceed 10MB.");
          postImageInput.value = ""; // Clear the input
          imagePreview.style.display = "none";
          return;
        }

        const reader = new FileReader();
        reader.onload = function (e) {
          imagePreview.src = e.target.result;
          imagePreview.style.display = "block"; // Show the image preview
        };
        reader.readAsDataURL(file);
      } else {
        imagePreview.src = "";
        imagePreview.style.display = "none";
      }
    });
  }
});

// Handle post creation
document.getElementById("newPostForm").addEventListener("submit", async (e) => {
  e.preventDefault();

  const formData = new FormData(e.target);
  const selectedCategories = Array.from(formData.getAll("categories")).map(
    Number
  );

  const postTitle = formData.get("title");
  const postContent = formData.get("content");
  const postImageFile = document.getElementById("postImageInput").files[0]; // Get the file directly

  const createPostErrorDiv = document.getElementById("create-post-error");
  createPostErrorDiv.style.display = "none";

  if (!postTitle || !postContent) {
    createPostErrorDiv.textContent = "Title and content are required.";
    createPostErrorDiv.style.display = "block";
    return;
  }

  let imageUrl = "";

  // If an image file is selected, upload it first
  if (postImageFile) {
    const imageFormData = new FormData();
    imageFormData.append("image", postImageFile);

    try {
      const uploadResponse = await fetch("/upload-image", {
        method: "POST",
        body: imageFormData,
      });

      if (!uploadResponse.ok) {
        const errorData = await uploadResponse.json();
        throw new Error(errorData.message || "Failed to upload image.");
      }

      const uploadResult = await uploadResponse.json();
      imageUrl = uploadResult.image_url;
      console.log("Image uploaded successfully:", imageUrl);
    } catch (error) {
      console.error("Image upload error:", error);
      createPostErrorDiv.textContent = `Image upload failed: ${error.message}`;
      createPostErrorDiv.style.display = "block";
      return; // Stop post creation if image upload fails
    }
  }

  // Now Create the Post plus the image URL
  const postData = {
    title: postTitle,
    content: postContent,
    categories: selectedCategories,
    image_url: imageUrl,
  };

  try {
    const response = await fetch("/posts", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(postData),
    });

    if (response.ok) {
      document.getElementById("create-post-form").style.display = "none";
      e.target.reset();
      document.getElementById("imagePreview").style.display = "none"; // Hide preview
      fetchPosts(); // Refresh posts lists
      alert("Post creation successful!");
    } else {
      const error = await response.text();
      document.getElementById("create-post-error").textContent = error;
      document.getElementById("create-post-error").style.display = "block";
    }
  } catch (error) {
    console.error("Error creating post:", error);
    document.getElementById("create-post-error").textContent =
      "Error creating post";
    document.getElementById("create-post-error").style.display = "block";
  }
});

// Like post
async function likePost(postId) {
  try {
    const response = await fetch("/posts/like", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ post_id: postId, is_like: true }),
    });

    if (response.ok) {
      fetchPosts();
    }
  } catch (error) {
    console.error("Error liking post:", error);
  }
}

// Dislike post
async function dislikePost(postId) {
  try {
    const response = await fetch("/posts/dislike", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ post_id: postId, is_like: false }),
    });

    if (response.ok) {
      fetchPosts(); // Refresh posts to update dislike count
    }
  } catch (error) {
    console.error("Error disliking post:", error);
  }
}

// Delete post
async function deletePost(postId) {
  if (!confirm("Are you sure you want to delete this post?")) {
    return;
  }

  try {
    const response = await fetch(`/posts?post_id=${postId}`, {
      method: "DELETE",
      credentials: "include",
    });

    if (response.ok) {
      fetchPosts();
    } else {
      const error = await response.text();
      alert(error || "Failed to delete post");
    }
  } catch (error) {
    console.error("Error deleting post:", error);
    alert("Error deleting post");
  }
}

// // Toggle comments section visibility
function toggleComments(postId) {
  const commentsSection = document.getElementById(`comments-section-${postId}`);
  if (!commentsSection) return;

  if (commentsSection.style.display === "none") {
    commentsSection.style.display = "block";
    initializeComments(postId);
  } else {
    commentsSection.style.display = "none";
  }
}

// Initialize the page
document.addEventListener("DOMContentLoaded", () => {});
