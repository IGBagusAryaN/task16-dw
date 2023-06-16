const promise = new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest();
    xhr.open("GET", "https://api.npoint.io/b88e52b88c0352deff42", true);
    xhr.onload = () => {
      if (xhr.status === 200) {
        resolve(JSON.parse(xhr.response));
      } else {
        reject("Error loading data.");
      }
    };
    xhr.onerror = () => {
      reject("Network error.");
    };
    xhr.send();
  });
  
  async function getAllTestimonials() {
    const response = await promise;
  
    let testimonialHTML = "";
    response.forEach(function (item) {
      testimonialHTML += `<div class="card col-lg-3 col-md-5 col-sm-5 p-3 m-2 rounded-3">
                              <img
                                src="${item.image}"
                                class="img-fluid profile-testimonial"
                              />
                              <p class="mt-2">${item.quote}</p>
                              <p class="text-end mb-1 fw-bold">${item.author}</p>
                              <p class="text-end fs-7">${item.rating} <i class="fa-solid fa-star fa-2xs"></i></p>
                          </div>
                          `;
    });
  
    document.getElementById("testimonials").innerHTML = testimonialHTML;
  }
  
  getAllTestimonials();
  
  async function getFilteredTestimonials(rating) {
    const response = await promise;
  
    const testimonialFiltered = response.filter((item) => {
      return item.rating === rating;
    });
  
  
    let testimonialHTML = "";
  
    if (testimonialFiltered.length === 0) {
      testimonialHTML = "<h1>Data not found!</h1>";
    } else {
      testimonialFiltered.forEach((item) => {
        testimonialHTML += `<div class="card col-lg-3 col-md-5 col-sm-5 p-3 m-2 rounded-3">
                                <img
                                  src="${item.image}"
                                  class="img-fluid profile-testimonial"
                                />
                                <p class="mt-2">${item.quote}</p>
                                <p class="text-end mb-1 fw-bold">${item.author}</p>
                                <p class="text-end fs-7">${item.rating} <i class="fa-solid fa-star fa-2xs"></i></p>
                            </div>
                          `;
      });
    }
  
    document.getElementById("testimonials").innerHTML = testimonialHTML;
  }