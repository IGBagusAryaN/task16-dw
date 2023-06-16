let dataBlog=[]

function addBlog(event){
    event.preventDefault()
        let title = document.getElementById('input-blog-title').value
        let startDate = new Date (document.getElementById('input-start-date').value)
        let endDate = new Date (document.getElementById('input-end-date').value)
        let description = document.getElementById('input-blog-desc').value
        let image = document.getElementById('image-upload').files

    let diff = Math.abs(endDate - startDate); 
    let months = Math.floor(diff / (1000 * 60 * 60 * 24 * 30)); 
    let days = Math.floor(diff / (1000 * 60 * 60 * 24)) % 30;


    
    const htmlIcon = '<img class="img-icon" src="./assets/images/svg/html5.svg">';
    const cssIcon = '<img class="img-icon" src="./assets/images/svg/css3.svg">';
    const reactjsIcon = '<img class="img-icon" src="./assets/images/svg/reactjs.svg">';
    const javascriptIcon = '<img class="img-icon" src="./assets/images/svg/javascript.svg">';

    let checkHtml = document.getElementById("input-checkbox-html").checked
        ? htmlIcon
        : "";
    let checkCss = document.getElementById("input-checkbox-css").checked
        ? cssIcon
        : "";
    let checkReact = document.getElementById("input-checkbox-rjs").checked
        ? reactjsIcon
        : "";
    let checkJs = document.getElementById("input-checkbox-js").checked
        ? javascriptIcon
        : "";


    image = URL.createObjectURL(image[0]);

    let blog = {
        title,
        months,
        days,
        description,
        image,
        checkHtml,
        checkCss,
        checkReact,
        checkJs,
        postAt : new Date(),

    }

    dataBlog.push(blog)
    console.log(dataBlog)

    renderBlog();
}

function renderBlog() {
    document.getElementById("contents").innerHTML = "";

    for(let index= 0; index < dataBlog.length; index++) {
        document.getElementById("contents").innerHTML += `
        
        <div class="card-project">
            <div class="img-project">
                <img src=${dataBlog[index].image} alt="" >
            </div>
            <a  href="./blog-detail/blog-detail1.html" class="title-project">${dataBlog[index].title}</a>
            <div class="duration-project">Duration : ${dataBlog[index].months} Bulan, ${dataBlog[index].days} Hari
            <div class="post-time-project">${getDistanceTime(dataBlog[index].postAt)} - ${getFullTime(dataBlog[index].postAt)}</div>
            <div class="mini-desc-project">${dataBlog[index].description}</div>
            <div class="tech-project">
                ${dataBlog[index].checkHtml}
                ${dataBlog[index].checkCss}
                ${dataBlog[index].checkReact}
                ${dataBlog[index].checkJs}
            </div>
            <div class="button-project">
                <button>Edit</button>
                <button>Delete</button>
            </div>
        </div>`
        
    }
}


function formSubmit(){
    let title = document.getElementById('input-blog-title').value
    let startDate = document.getElementById('input-start-date').value
    let endDate = document.getElementById('input-end-date').value
    let description = document.getElementById('input-blog-desc').value
    let image = document.getElementById('image-upload').files
  
    if (title == "") {
      return alert("Name Project Required!");
    } else if (startDate == "") {
      return alert("Start date Required!");
    } else if (endDate == "") {
      return alert("End date Required!");
    } else if (description == "") {
      return alert("Description Required!");
    } else if (image == "") {
      return alert("Attach File Required");
    }
}



function getFullTime(time) {
    let monthName = [
      "Jan",
      "Feb",
      "Mar",
      "Apr",
      "May",
      "Jun",
      "Jul",
      "Aug",
      "Sep",
      "Oct",
      "Nov",
      "Dec",
    ];
  
    let date = time.getDate();
    let monthIndex = time.getMonth();
    let year = time.getFullYear();
    let hours = time.getHours();
    let minutes = time.getMinutes();

  
    if (hours <= 9) {
      hours = "0" + hours;
    } else if (minutes <= 9) {
      minutes = "0" + minutes;
    }
  
    return `${date} ${monthName[monthIndex]} ${year} ${hours}:${minutes} WITA`;
  }
  
  function getDistanceTime(time) {
    let timeNow = new Date();
    let timePost = time;
  
    let distance = timeNow - timePost; 
    console.log(distance);
  
    let milisecond = 1000; 
    let secondInHours = 3600; 
    let hoursInDays = 24; 
  
    let distanceDay = Math.floor(
      distance / (milisecond * secondInHours * hoursInDays)
    );
    let distanceHours = Math.floor(distance / (milisecond * 60 * 60)); 
    let distanceMinutes = Math.floor(distance / (milisecond * 60)); 
    let distanceSeconds = Math.floor(distance / milisecond); 
  
    if (distanceDay > 0) {
      return `${distanceDay} Day Ago`;
    } else if (distanceHours > 0) {
      return `${distanceHours} Hours Ago`;
    } else if (distanceMinutes > 0) {
      return `${distanceMinutes} Minutes Ago`;
    } else {
      return `${distanceSeconds} Seconds Ago`;
    }

  
  }
  // setInterval(function () {
  //   renderBlog();
  // }, 10000);
 