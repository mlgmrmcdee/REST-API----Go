async function loadMovie() {
  const params = new URLSearchParams(window.location.search);
  const id = params.get("id");
  if (!id) return;

  const response = await fetch(`/movie?id=${id}`);
  const movie = await response.json();

  document.getElementById("poster").src =
    movie.poster || "https://via.placeholder.com/300x450";

  document.getElementById("title").textContent = movie.title;
  document.getElementById("description").textContent =
    movie.description || "No description";

  document.getElementById("year").textContent = movie.year || "-";
  document.getElementById("duration").textContent = movie.duration_min || "-";

  document.getElementById("rating").textContent = movie.rating || "-";
  document.getElementById("certificate").textContent = movie.certificate || "-";

  document.getElementById("genres").textContent =
    movie.genres?.join(", ") || "-";

  document.getElementById("directors").textContent =
    movie.directors?.join(", ") || "-";

  document.getElementById("actors").textContent =
    movie.actors?.join(", ") || "-";
}

loadMovie();
