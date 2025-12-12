document.addEventListener("DOMContentLoaded", async () => {
  const renderedIds = new Set();

  let currentPage = 1;
  const limit = 25;
  let isLoading = false;

  const filters = new URLSearchParams();

  // Элементы фильтров
  const genreEl = document.getElementById("genre");
  const certificateEl = document.getElementById("certificate");
  const yearFromEl = document.getElementById("year_from");
  const yearToEl = document.getElementById("year_to");
  const sortEl = document.getElementById("sort");
  const applyBtn = document.getElementById("applyFilters");

  if (
    !genreEl ||
    !certificateEl ||
    !yearFromEl ||
    !yearToEl ||
    !sortEl ||
    !applyBtn
  ) {
    console.error("Один или несколько элементов фильтров не найдены в DOM");
    return;
  }

  // Подгружаем жанры
  async function loadGenres() {
    try {
      const res = await fetch("/genres");
      if (!res.ok) throw new Error("Ошибка при загрузке жанров");
      const genres = await res.json();

      genres.forEach((g) => {
        const option = document.createElement("option");
        option.value = g;
        option.textContent = g;
        genreEl.appendChild(option);
      });
    } catch (err) {
      console.error(err);
    }
  }

  async function loadCertificates() {
    try {
      const res = await fetch("/certificates");
      if (!res.ok) throw new Error("Ошибка при загрузке возрастных рейтингов");
      const certs = await res.json();

      certs.forEach((c) => {
        const option = document.createElement("option");
        option.value = c;
        option.textContent = c;
        certificateEl.appendChild(option);
      });
    } catch (err) {
      console.error(err);
    }
  }

  // Получение фильмов
  async function fetchMovies(page, reset = false) {
    if (isLoading) return;
    isLoading = true;

    try {
      const container = document.getElementById("movies");
      if (reset) container.innerHTML = "";

      let moviesWithPoster = [];
      let pageToFetch = page;

      while (moviesWithPoster.length < limit) {
        filters.set("page", pageToFetch);
        filters.set("limit", limit);

        const res = await fetch(`/movies/filter?${filters.toString()}`);
        if (!res.ok) throw new Error("Ошибка при загрузке фильмов");

        const data = await res.json();
        if (data.length === 0) break;

        for (const movie of data) {
          if (moviesWithPoster.length >= limit) break;
          if (!movie.poster) continue;

          const exists = await new Promise((resolve) => {
            const img = new Image();
            img.src = movie.poster;
            img.onload = () => resolve(true);
            img.onerror = () => resolve(false);
          });

          if (exists) moviesWithPoster.push(movie);
        }

        pageToFetch++;
      }

      renderMovies(moviesWithPoster, !reset);
      currentPage = pageToFetch;
    } catch (err) {
      const container = document.getElementById("movies");
      if (container.children.length === 0) {
        container.innerHTML = "<p>Не удалось загрузить фильмы.</p>";
      }
    } finally {
      isLoading = false;
    }
  }

  // Рендер карточек
  function renderMovies(movies, append = false) {
    const container = document.getElementById("movies");
    if (!append) {
      container.innerHTML = "";
      renderedIds.clear(); // очищаем трекер если reset
    }

    movies.forEach((movie) => {
      if (renderedIds.has(movie.id)) return; // <<< не рендерим дубликаты
      renderedIds.add(movie.id);

      const card = document.createElement("div");
      card.className = "movie-card";
      card.onclick = () => {
        window.location.href = `/static/movie.html?id=${movie.id}`;
      };

      const img = document.createElement("img");
      img.className = "movie-card__img";
      img.src = movie.poster;
      card.appendChild(img);

      const body = document.createElement("div");
      body.className = "movie-card__body";

      const title = document.createElement("div");
      title.className = "movie-card__title";
      title.textContent = movie.title;
      body.appendChild(title);

      const info = document.createElement("div");
      info.className = "movie-card__info";
      info.textContent = `${movie.year || "N/A"} | ${
        movie.duration_min || "N/A"
      } min | ${movie.certificate || "N/A"}`;
      body.appendChild(info);

      const rating = document.createElement("div");
      rating.className = "movie-card__rating";
      rating.innerHTML = `<span class="star">★</span> ${movie.rating || "N/A"}`;
      body.appendChild(rating);

      card.appendChild(body);
      container.appendChild(card);
    });
  }

  // Кнопка Apply
  applyBtn.addEventListener("click", () => {
    filters.delete("genre");
    filters.delete("certificate");
    filters.delete("year_from");
    filters.delete("year_to");
    filters.delete("sort");

    if (genreEl.value) filters.set("genre", genreEl.value);
    if (certificateEl.value) filters.set("certificate", certificateEl.value);
    if (yearFromEl.value) filters.set("year_from", yearFromEl.value);
    if (yearToEl.value) filters.set("year_to", yearToEl.value);
    if (sortEl.value) filters.set("sort", sortEl.value);

    currentPage = 1;
    fetchMovies(currentPage, true);
  });

  // Бесконечная прокрутка
  window.addEventListener("scroll", () => {
    const scrollTop = window.scrollY || document.documentElement.scrollTop;
    const windowHeight = window.innerHeight;
    const fullHeight = document.documentElement.scrollHeight;

    if (scrollTop + windowHeight >= fullHeight - 5) {
      // 5px запас
      fetchMovies(currentPage);
    }
  });

  // Инициализация
  await loadGenres();
  await loadCertificates();
  fetchMovies(currentPage, true);
});
