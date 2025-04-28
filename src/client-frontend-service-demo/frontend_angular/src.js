// API URLs
const API_BASE_URL = "http://localhost:8080"; // In a real scenario, this would be properly configured
const API_ENDPOINTS = {
  users: `${API_BASE_URL}/api/users`,
  restaurants: "http://localhost:8081/api/restaurants",
  orders: "http://localhost:8082/api/orders",
  payments: "http://localhost:8083/api/payments",
  deliveries: "http://localhost:8084/api/deliveries",
  notifications: "http://localhost:8085/api/notifications",
};

// Global variables
let currentUser = null;
let currentRestaurant = null;
let currentOrder = null;
let cart = {
  items: [],
  restaurantId: null,
};

// On DOM load
document.addEventListener("DOMContentLoaded", function () {
  // Load logged in user (assume user ID 1 for demo)
  fetchUser(1);

  // Initialize navigation
  initNavigation();

  // Load restaurants
  fetchRestaurants();

  // Initialize checkout button
  document
    .getElementById("checkout-btn")
    .addEventListener("click", function () {
      prepareCheckout();
      showSection("checkout");
    });

  // Initialize place order button
  document
    .getElementById("place-order-btn")
    .addEventListener("click", placeOrder);

  // Initialize search functionality
  document.querySelector(".search-btn").addEventListener("click", function () {
    const searchQuery = document
      .getElementById("restaurant-search")
      .value.trim();
    if (searchQuery) {
      // In a real app, we would search from the API
      // For this demo, we just log the search query
      console.log(`Searching for: ${searchQuery}`);
    }
  });
});

// Initialize navigation
function initNavigation() {
  const navLinks = document.querySelectorAll("nav a");

  navLinks.forEach((link) => {
    link.addEventListener("click", function (e) {
      e.preventDefault();

      // Remove active class from all links
      navLinks.forEach((l) => l.classList.remove("active"));

      // Add active class to the clicked link
      this.classList.add("active");

      // Show the corresponding section
      const sectionToShow = this.getAttribute("data-section");
      showSection(sectionToShow);

      // Load section data if needed
      if (sectionToShow === "restaurants") {
        fetchRestaurants();
      } else if (sectionToShow === "orders") {
        fetchOrders(currentUser.id);
      } else if (sectionToShow === "profile") {
        displayProfile();
      }
    });
  });
}

// Show a specific section and hide others
function showSection(sectionId) {
  const sections = document.querySelectorAll("main section");

  sections.forEach((section) => {
    if (section.id === sectionId) {
      section.classList.remove("hidden-section");
      section.classList.add("active-section");
    } else {
      section.classList.remove("active-section");
      section.classList.add("hidden-section");
    }
  });
}

// Fetch user data
async function fetchUser(userId) {
  try {
    const response = await fetch(`${API_ENDPOINTS.users}/${userId}`);

    if (!response.ok) {
      throw new Error("Failed to fetch user data");
    }

    const userData = await response.json();
    currentUser = userData;

    // For demo purposes, if we can't connect to the API, use a mock user
    if (!currentUser) {
      currentUser = {
        id: 1,
        name: "John Doe",
        email: "john.doe@example.com",
        phone: "555-1234",
        address: "123 Main St, City",
      };
    }
  } catch (error) {
    console.error("Error fetching user:", error);

    // Use mock data for demo
    currentUser = {
      id: 1,
      name: "John Doe",
      email: "john.doe@example.com",
      phone: "555-1234",
      address: "123 Main St, City",
    };
  }
}

// Fetch restaurants
async function fetchRestaurants() {
  const restaurantsContainer = document.getElementById("restaurants-container");
  restaurantsContainer.innerHTML =
    '<div class="loading">Loading restaurants...</div>';

  try {
    const response = await fetch(API_ENDPOINTS.restaurants);

    if (!response.ok) {
      throw new Error("Failed to fetch restaurants");
    }

    const restaurants = await response.json();
    displayRestaurants(restaurants);
  } catch (error) {
    console.error("Error fetching restaurants:", error);

    // // Use mock data for demo
    // const mockRestaurants = [
    //     {
    //         id: 1,
    //         name: 'Tasty Bites',
    //         address: '123 Main St',
    //         cuisine: 'Italian',
    //         rating: 4.5
    //     },
    //     {
    //         id: 2,
    //         name: 'Burger Palace',
    //         address: '456 Oak Ave',
    //         cuisine: 'American',
    //         rating: 4.2
    //     },
    //     {
    //         id: 3,
    //         name: 'Sushi Heaven',
    //         address: '789 Maple Rd',
    //         cuisine: 'Japanese',
    //         rating: 4.7
    //     }
    // ];

    // displayRestaurants(mockRestaurants);
  }
}

// Display restaurants
function displayRestaurants(restaurants) {
  const restaurantsContainer = document.getElementById("restaurants-container");

  if (!restaurants || restaurants.length === 0) {
    restaurantsContainer.innerHTML = "<p>No restaurants found.</p>";
    return;
  }

  let html = "";

  restaurants.forEach((restaurant) => {
    html += `
            <div class="restaurant-card" onclick="fetchRestaurantDetails(${
              restaurant.id
            })">
                <div class="restaurant-image">Restaurant Image</div>
                <div class="restaurant-info">
                    <div class="restaurant-name">${restaurant.name}</div>
                    <div class="restaurant-cuisine">${restaurant.cuisine}</div>
                    <div class="restaurant-rating">
                        <span class="rating">${restaurant.rating.toFixed(
                          1
                        )}</span>
                        <span>${restaurant.address}</span>
                    </div>
                </div>
            </div>
        `;
  });

  restaurantsContainer.innerHTML = html;
}

// Fetch restaurant details
async function fetchRestaurantDetails(restaurantId) {
  const restaurantDetailContainer = document.getElementById(
    "restaurant-detail-container"
  );
  const menuItemsContainer = document.getElementById("menu-items-container");

  restaurantDetailContainer.innerHTML =
    '<div class="loading">Loading restaurant details...</div>';
  menuItemsContainer.innerHTML = '<div class="loading">Loading menu...</div>';

  try {
    // Fetch restaurant details
    const restaurantResponse = await fetch(
      `${API_ENDPOINTS.restaurants}/${restaurantId}`
    );

    if (!restaurantResponse.ok) {
      throw new Error("Failed to fetch restaurant details");
    }

    const restaurant = await restaurantResponse.json();
    currentRestaurant = restaurant;

    // Display restaurant details
    displayRestaurantDetails(restaurant);

    // Display menu items
    displayMenuItems(restaurant.menuItems);

    // Show restaurant detail section
    showSection("restaurant-detail");

    // Reset cart if it's a different restaurant
    if (cart.restaurantId !== restaurantId) {
      cart.items = [];
      cart.restaurantId = restaurantId;
      updateCart();
    }
  } catch (error) {
    console.error("Error fetching restaurant details:", error);

    // Use mock data for demo
    const mockRestaurant = {
      id: restaurantId,
      name: "Tasty Bites",
      address: "123 Main St",
      cuisine: "Italian",
      rating: 4.5,
      menuItems: [
        {
          id: 1,
          name: "Margherita Pizza",
          description: "Classic pizza with tomato sauce, mozzarella, and basil",
          price: 12.99,
          category: "Main",
        },
        {
          id: 2,
          name: "Spaghetti Carbonara",
          description: "Pasta with egg, cheese, pancetta, and black pepper",
          price: 14.99,
          category: "Main",
        },
        {
          id: 3,
          name: "Tiramisu",
          description: "Coffee-flavored Italian dessert",
          price: 7.99,
          category: "Dessert",
        },
        {
          id: 4,
          name: "Caesar Salad",
          description:
            "Romaine lettuce, croutons, parmesan cheese, and Caesar dressing",
          price: 8.99,
          category: "Starter",
        },
      ],
    };

    currentRestaurant = mockRestaurant;

    // Display restaurant details
    displayRestaurantDetails(mockRestaurant);

    // Display menu items
    displayMenuItems(mockRestaurant.menuItems);

    // Show restaurant detail section
    showSection("restaurant-detail");

    // Reset cart if it's a different restaurant
    if (cart.restaurantId !== restaurantId) {
      cart.items = [];
      cart.restaurantId = restaurantId;
      updateCart();
    }
  }
}

// Display restaurant details
function displayRestaurantDetails(restaurant) {
  const restaurantDetailContainer = document.getElementById(
    "restaurant-detail-container"
  );

  const html = `
        <div class="restaurant-detail">
            <div class="restaurant-header">
                <div>
                    <h2>${restaurant.name}</h2>
                    <p>${restaurant.cuisine} • ${restaurant.address}</p>
                </div>
                <div class="restaurant-actions">
                    <span class="rating">${restaurant.rating.toFixed(1)}</span>
                </div>
            </div>
        </div>
    `;

  restaurantDetailContainer.innerHTML = html;
}

// Display menu items
function displayMenuItems(menuItems) {
  const menuItemsContainer = document.getElementById("menu-items-container");

  if (!menuItems || menuItems.length === 0) {
    menuItemsContainer.innerHTML = "<p>No menu items available.</p>";
    return;
  }

  // Group menu items by category
  const menuItemsByCategory = menuItems.reduce((acc, item) => {
    const category = item.category || "Other";
    if (!acc[category]) {
      acc[category] = [];
    }
    acc[category].push(item);
    return acc;
  }, {});

  let html = "";

  // Display each category and its items
  Object.keys(menuItemsByCategory).forEach((category) => {
    html += `
            <div class="menu-section">
                <div class="menu-category">${category}</div>
                <div class="menu-items">
        `;

    menuItemsByCategory[category].forEach((item) => {
      html += `
                <div class="menu-item">
                    <div class="menu-item-name">${item.name}</div>
                    <div class="menu-item-description">${item.description}</div>
                    <div class="menu-item-price">$${item.price.toFixed(2)}</div>
                    <button class="add-to-cart" onclick="addToCart(${
                      item.id
                    }, '${item.name}', ${item.price})">+</button>
                </div>
            `;
    });

    html += `
                </div>
            </div>
        `;
  });

  menuItemsContainer.innerHTML = html;
}

// Add item to cart
function addToCart(itemId, itemName, itemPrice) {
  // Check if item already in cart
  const existingItemIndex = cart.items.findIndex((item) => item.id === itemId);

  if (existingItemIndex >= 0) {
    // Increment quantity if item already exists
    cart.items[existingItemIndex].quantity++;
  } else {
    // Add new item to cart
    cart.items.push({
      id: itemId,
      name: itemName,
      price: itemPrice,
      quantity: 1,
    });
  }

  // Update cart display
  updateCart();
}

// Update cart display
function updateCart() {
  const cartItemsContainer = document.getElementById("cart-items");
  const cartTotalElement = document.getElementById("cart-total-price");
  const checkoutBtn = document.getElementById("checkout-btn");

  if (cart.items.length === 0) {
    cartItemsContainer.innerHTML =
      '<p class="empty-cart">Your cart is empty</p>';
    cartTotalElement.textContent = "$0.00";
    checkoutBtn.disabled = true;
    return;
  }

  let html = "";
  let total = 0;

  cart.items.forEach((item) => {
    const itemTotal = item.price * item.quantity;
    total += itemTotal;

    html += `
            <div class="cart-item">
                <div class="cart-item-details">
                    <div class="cart-item-name">${item.name}</div>
                    <div class="cart-item-price">$${item.price.toFixed(2)}</div>
                </div>
                <div class="cart-item-actions">
                    <button class="quantity-btn" onclick="updateCartItemQuantity(${
                      item.id
                    }, ${item.quantity - 1})">-</button>
                    <div class="cart-item-quantity">${item.quantity}</div>
                    <button class="quantity-btn" onclick="updateCartItemQuantity(${
                      item.id
                    }, ${item.quantity + 1})">+</button>
                </div>
            </div>
        `;
  });

  cartItemsContainer.innerHTML = html;
  cartTotalElement.textContent = `$${total.toFixed(2)}`;
  checkoutBtn.disabled = false;
}

// Update cart item quantity
function updateCartItemQuantity(itemId, newQuantity) {
  if (newQuantity <= 0) {
    // Remove item from cart if quantity is 0 or less
    cart.items = cart.items.filter((item) => item.id !== itemId);
  } else {
    // Update quantity
    const item = cart.items.find((item) => item.id === itemId);
    if (item) {
      item.quantity = newQuantity;
    }
  }

  // Update cart display
  updateCart();
}

// Prepare checkout
function prepareCheckout() {
  const checkoutItemsContainer = document.getElementById("checkout-items");
  const checkoutTotalElement = document.getElementById("checkout-total-price");

  let html = "";
  let total = 0;

  cart.items.forEach((item) => {
    const itemTotal = item.price * item.quantity;
    total += itemTotal;

    html += `
            <div class="cart-item">
                <div class="cart-item-details">
                    <div class="cart-item-name">${item.name} x ${
      item.quantity
    }</div>
                </div>
                <div>$${itemTotal.toFixed(2)}</div>
            </div>
        `;
  });

  checkoutItemsContainer.innerHTML = html;
  checkoutTotalElement.textContent = `$${total.toFixed(2)}`;

  // Pre-fill address with user's address
  document.getElementById("delivery-address").value = currentUser.address;
}

// Place order
async function placeOrder() {
  // Get delivery address
  const deliveryAddress = document
    .getElementById("delivery-address")
    .value.trim();

  if (!deliveryAddress) {
    alert("Please enter a delivery address");
    return;
  }

  // Get payment method
  const paymentMethod = document.querySelector(
    'input[name="payment-method"]:checked'
  ).value;

  // Create order items
  const orderItems = cart.items.map((item) => {
    return {
      menuItemId: item.id,
      name: item.name,
      price: item.price,
      quantity: item.quantity,
    };
  });

  // Calculate total amount
  const totalAmount = cart.items.reduce(
    (total, item) => total + item.price * item.quantity,
    0
  );

  // Create order object
  const order = {
    userId: currentUser.id,
    restaurantId: cart.restaurantId,
    items: orderItems,
    totalAmount: totalAmount,
    status: "created",
    address: deliveryAddress,
  };

  try {
    // Create order
    const orderResponse = await fetch(API_ENDPOINTS.orders, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(order),
    });

    if (!orderResponse.ok) {
      throw new Error("Failed to create order");
    }

    const createdOrder = await orderResponse.json();
    currentOrder = createdOrder;

    // Display order confirmation
    displayOrderConfirmation(createdOrder);

    // Clear cart
    cart.items = [];
    cart.restaurantId = null;

    // Show confirmation section
    showSection("order-confirmation");
  } catch (error) {
    console.error("Error creating order:", error);

    // Mock order creation for demo
    const mockCreatedOrder = {
      id: Math.floor(Math.random() * 1000),
      userId: currentUser.id,
      restaurantId: cart.restaurantId,
      items: orderItems,
      totalAmount: totalAmount,
      status: "created",
      address: deliveryAddress,
      createdAt: new Date().toISOString(),
    };

    currentOrder = mockCreatedOrder;

    // Display order confirmation
    displayOrderConfirmation(mockCreatedOrder);

    // Clear cart
    cart.items = [];
    cart.restaurantId = null;

    // Show confirmation section
    showSection("order-confirmation");
  }
}

// Display order confirmation
function displayOrderConfirmation(order) {
  const confirmationDetailsContainer = document.getElementById(
    "order-confirmation-details"
  );

  let orderItemsHtml = "";
  order.items.forEach((item) => {
    orderItemsHtml += `
            <div class="confirmation-item">
                <span>${item.name} x ${item.quantity}</span>
                <span>$${(item.price * item.quantity).toFixed(2)}</span>
            </div>
        `;
  });

  const html = `
        <div class="confirmation-item">
            <span><strong>Order ID:</strong></span>
            <span>#${order.id}</span>
        </div>
        <div class="confirmation-item">
            <span><strong>Restaurant:</strong></span>
            <span>${currentRestaurant.name}</span>
        </div>
        <div class="confirmation-item">
            <span><strong>Delivery Address:</strong></span>
            <span>${order.address}</span>
        </div>
        <div class="confirmation-item">
            <span><strong>Order Date:</strong></span>
            <span>${new Date(order.createdAt).toLocaleString()}</span>
        </div>
        <h4 style="margin-top: 15px; margin-bottom: 10px;">Order Items</h4>
        ${orderItemsHtml}
        <div class="confirmation-item" style="font-weight: bold; margin-top: 10px;">
            <span>Total Amount:</span>
            <span>$${order.totalAmount.toFixed(2)}</span>
        </div>
    `;

  confirmationDetailsContainer.innerHTML = html;
}

// Fetch orders for user
async function fetchOrders(userId) {
  const ordersContainer = document.getElementById("orders-container");
  ordersContainer.innerHTML = '<div class="loading">Loading orders...</div>';

  try {
    const response = await fetch(`${API_ENDPOINTS.users}/${userId}/orders`);

    if (!response.ok) {
      throw new Error("Failed to fetch orders");
    }

    const orders = await response.json();
    displayOrders(orders);
  } catch (error) {
    console.error("Error fetching orders:", error);

    // Use mock data for demo
    const mockOrders = [
      {
        id: 1001,
        userId: userId,
        restaurantId: 1,
        items: [
          {
            menuItemId: 1,
            name: "Margherita Pizza",
            price: 12.99,
            quantity: 2,
          },
        ],
        totalAmount: 25.98,
        status: "delivered",
        address: "123 Main St, City",
        createdAt: "2025-04-24T14:30:00Z",
      },
      {
        id: 1002,
        userId: userId,
        restaurantId: 2,
        items: [
          {
            menuItemId: 5,
            name: "Cheeseburger",
            price: 9.99,
            quantity: 1,
          },
          {
            menuItemId: 6,
            name: "French Fries",
            price: 3.99,
            quantity: 1,
          },
        ],
        totalAmount: 13.98,
        status: "out_for_delivery",
        address: "123 Main St, City",
        createdAt: "2025-04-26T18:15:00Z",
      },
    ];

    displayOrders(mockOrders);
  }
}

// Display orders
function displayOrders(orders) {
  const ordersContainer = document.getElementById("orders-container");

  if (!orders || orders.length === 0) {
    ordersContainer.innerHTML = "<p>No orders found.</p>";
    return;
  }

  let html = "";

  // Sort orders by date (newest first)
  orders.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));

  orders.forEach((order) => {
    html += `
            <div class="order-card" onclick="fetchOrderDetails(${order.id})">
                <div class="order-header">
                    <div class="order-number">#${order.id}</div>
                    <div class="order-date">${new Date(
                      order.createdAt
                    ).toLocaleDateString()}</div>
                </div>
                <div class="order-status status-${
                  order.status
                }">${formatOrderStatus(order.status)}</div>
                <div class="order-restaurant">Restaurant: ${getRestaurantName(
                  order.restaurantId
                )}</div>
                <div class="order-total">$${order.totalAmount.toFixed(2)}</div>
            </div>
        `;
  });

  ordersContainer.innerHTML = html;
}

// Format order status for display
function formatOrderStatus(status) {
  switch (status) {
    case "created":
      return "Created";
    case "paid":
      return "Paid";
    case "preparing":
      return "Preparing";
    case "out_for_delivery":
      return "Out for Delivery";
    case "delivered":
      return "Delivered";
    case "cancelled":
      return "Cancelled";
    default:
      return status.charAt(0).toUpperCase() + status.slice(1);
  }
}

// Get restaurant name by ID (mock function)
function getRestaurantName(restaurantId) {
  // In a real app, this would fetch from cache or API
  const mockRestaurantNames = {
    1: "Tasty Bites",
    2: "Burger Palace",
    3: "Sushi Heaven",
  };

  return mockRestaurantNames[restaurantId] || "Unknown Restaurant";
}

// Fetch order details
async function fetchOrderDetails(orderId) {
  const orderDetailContainer = document.getElementById(
    "order-detail-container"
  );
  orderDetailContainer.innerHTML =
    '<div class="loading">Loading order details...</div>';

  try {
    const response = await fetch(`${API_ENDPOINTS.orders}/${orderId}`);

    if (!response.ok) {
      throw new Error("Failed to fetch order details");
    }

    const order = await response.json();
    displayOrderDetails(order);

    // Show order detail section
    showSection("order-detail");
  } catch (error) {
    console.error("Error fetching order details:", error);

    // Use mock data for demo
    const mockOrder = {
      id: orderId,
      userId: currentUser.id,
      restaurantId: 1,
      items: [
        {
          menuItemId: 1,
          name: "Margherita Pizza",
          price: 12.99,
          quantity: 2,
        },
      ],
      totalAmount: 25.98,
      status: "delivered",
      address: "123 Main St, City",
      createdAt: "2025-04-24T14:30:00Z",
      updatedAt: "2025-04-24T15:45:00Z",
    };

    displayOrderDetails(mockOrder);

    // Show order detail section
    showSection("order-detail");
  }
}

// Display order details
function displayOrderDetails(order) {
  const orderDetailContainer = document.getElementById(
    "order-detail-container"
  );

  let itemsHtml = "";
  order.items.forEach((item) => {
    itemsHtml += `
            <div class="cart-item">
                <div class="cart-item-details">
                    <div class="cart-item-name">${item.name}</div>
                    <div class="cart-item-price">$${item.price.toFixed(2)} x ${
      item.quantity
    }</div>
                </div>
                <div>$${(item.price * item.quantity).toFixed(2)}</div>
            </div>
        `;
  });

  const html = `
        <div class="restaurant-detail">
            <h2>Order #${order.id}</h2>
            <div class="order-status status-${
              order.status
            }" style="margin: 10px 0;">${formatOrderStatus(order.status)}</div>
            
            <div style="margin-bottom: 20px;">
                <div><strong>Restaurant:</strong> ${getRestaurantName(
                  order.restaurantId
                )}</div>
                <div><strong>Date:</strong> ${new Date(
                  order.createdAt
                ).toLocaleString()}</div>
                <div><strong>Delivery Address:</strong> ${order.address}</div>
            </div>
            
            <h3>Order Items</h3>
            <div style="margin-top: 10px;">
                ${itemsHtml}
            </div>
            
            <div class="cart-total" style="margin-top: 20px;">
                <span>Total Amount:</span>
                <span>$${order.totalAmount.toFixed(2)}</span>
            </div>
        </div>
    `;

  orderDetailContainer.innerHTML = html;
}

// Display user profile
function displayProfile() {
  const profileContainer = document.getElementById("profile-container");

  if (!currentUser) {
    profileContainer.innerHTML =
      '<div class="loading">Loading profile...</div>';
    return;
  }

  const html = `
        <div class="profile-section">
            <h3>Personal Information</h3>
            <div class="profile-info">
                <div class="profile-label">Name</div>
                <div class="profile-value">${currentUser.name}</div>
            </div>
            <div class="profile-info">
                <div class="profile-label">Email</div>
                <div class="profile-value">${currentUser.email}</div>
            </div>
            <div class="profile-info">
                <div class="profile-label">Phone</div>
                <div class="profile-value">${currentUser.phone}</div>
            </div>
            <div class="profile-info">
                <div class="profile-label">Address</div>
                <div class="profile-value">${currentUser.address}</div>
            </div>
            <button class="edit-profile">Edit Profile</button>
        </div>
    `;

  profileContainer.innerHTML = html;
}
