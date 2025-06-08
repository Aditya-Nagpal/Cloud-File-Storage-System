<template>
  <nav class="navbar navbar-light bg-white border-bottom px-4 d-flex justify-content-between">
    <h4 class="mt-2 fw-bold">FastFiles</h4>

    <div class="d-flex align-items-center gap-3">
      <!-- Upload dropdown -->
      <div class="dropdown">
        <button
          class="btn btn-primary dropdown-toggle"
          type="button"
          data-bs-toggle="dropdown"
        >
          Upload
        </button>
        <ul class="dropdown-menu">
          <li><a class="dropdown-item" href="#">Upload File</a></li>
          <li><a class="dropdown-item" href="#">Upload Folder</a></li>
        </ul>
      </div>

      <!-- Profile dropdown -->
      <div class="dropdown">
        <img
          src="https://via.placeholder.com/40"
          class="rounded-circle"
          style="width: 40px; height: 40px; cursor: pointer;"
          data-bs-toggle="dropdown"
        />
        <ul class="dropdown-menu dropdown-menu-end">
          <li><a class="dropdown-item" href="#">View Profile</a></li>
          <li><a class="dropdown-item" href="#">Logout</a></li>
        </ul>
      </div>
    </div>
  </nav>
</template>
