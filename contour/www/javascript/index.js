window.addEventListener("load", function load(event){
    
    const upload_el = document.getElementById("upload");
    const feedback_el = document.getElementById("feedback");
    const results_el = document.getElementById("results");    
    
    const image_btn = document.getElementById("image-button");
    const start_video_btn = document.getElementById("start-video");
    const stop_video_btn = document.getElementById("stop-video");
    const contour_video_btn = document.getElementById("contour-video");                

    const iterations_el = document.getElementById("iterations");
    const iterations_count_el = document.getElementById("iterations-count");    
    const video = document.getElementById("video");

    const image_spinner = document.getElementById("contour-image-spinner-svg");
    const video_spinner = document.getElementById("contour-video-spinner-svg");    
    
    var video_b64;

    const parser = new DOMParser();
    
    sfomuseum.golang.wasm.fetch("wasm/contour.wasm").then((rsp) => {

	iterations_el.onchange = function(e){
	    const el = e.target;
	    iterations_count_el.innerText = el.value;
	};

	iterations_count_el.innerText = iterations_el.value;
	
	var do_contour_image = function(im_b64, n){

	    return new Promise((resolve, reject) => {
		
		results_el.innerHTML = "";
		console.debug("Start contour");

		contour_svg(im_b64, n).then((rsp) => {
		    
		    console.debug("Contour as SVG successful", rsp);
		    
		    const doc = parser.parseFromString(rsp, "image/svg+xml");
		    results_el.appendChild(doc.documentElement);

		    resolve();
		    
		}).catch((err) => {
		    feedback_el.innerText = "Failed to contour image, " + err;
		    console.error("Failed to contour as SVG", err);
		    reject(err);
		});
		
	    });
	};
	
	var process_video_tick = function(){

	    if (video.readyState === video.HAVE_ENOUGH_DATA) {

		const canvas = document.createElement("canvas");
		const context = canvas.getContext('2d');
		
		canvas.width = video.videoWidth;
		canvas.height = video.videoHeight;
		
		context.drawImage(video, 0, 0, canvas.width, canvas.height);
		video_b64 = canvas.toDataURL('image/png').replace("data:image/png;base64,", "");		
	    }
	    
	    requestAnimationFrame(process_video_tick);
	};
	
	var process_video = function(stream){

	    video.style.display = "block";
	    
	    video.srcObject = stream;
	    video.setAttribute("playsinline", true); // required to tell iOS safari we don't want fullscreen
	    video.play();
	    
	    requestAnimationFrame(process_video_tick);
	}
	
	var process_upload = function(){

	    if (! upload_el.files.length){
		feedback_el.innerText = "There are no files to process";
		return;
	    }
	    
	    const file = upload_el.files[0];
	    
	    if (! file.type.startsWith('image/')){
		return false;
	    }

	    switch (file.type) {
		case "image/jpeg":
		case "image/png":
		case "image/gif":
		case "image/webp":
		    // pass
		    break;
		default:
		    feedback_el.innerText = "Unsupported file type: " + file.type;
		    return;
	    }
	    
            const reader = new FileReader();

            reader.onload = function(e) {
		
		const img = document.createElement("img");
		img.setAttribute("style", "max-height: 400px; max-width:400px;");
		img.src = e.target.result;
		
		const wrapper = document.getElementById("image-wrapper");
		wrapper.innerHTML = "";
		wrapper.appendChild(img);
            };
	    
            reader.readAsDataURL(file);

	    setTimeout(function(){

		reader.onload = function(e) {
		    const im_b64 = e.target.result;
		    const prefix = "data:" + file.type + ";base64,";
		    const iterations = iterations_el.valueAsNumber;

		    image_spinner.style.display = "inline-block";
		    image_btn.setAttribute("disabled", "disabled");
		    
		    setTimeout(function(){
			
			do_contour_image(im_b64.replace(prefix, ""), iterations).then((rsp) => {
			    image_spinner.style.display = "none";
			    image_btn.removeAttribute("disabled");
			}).catch((err) => {
			    image_spinner.style.display = "none";
			    image_btn.removeAttribute("disabled");

			    feedback_el.innerText = "Failed to contour image, " + err;			    
			});
			
		    }, 10);
		};

		reader.readAsDataURL(file);
		
	    }, 10)
	    
	};

	upload_el.onchange = function(){
	    results_el.innerHTML = "";	    
	};

	contour_video_btn.onclick = function(){
	    
	    const iterations = iterations_el.valueAsNumber;

	    video_spinner.style.display = "inline-block";
	    contour_video_btn.setAttribute("disbled", "disabled");
	    
	    setTimeout(function(){
		
		do_contour_image(video_b64, iterations).then((rsp) => {
		    video_spinner.style.display = "none";
		    contour_video_btn.removeAttribute("disabled");
		}).catch((err) => {
		    video_spinner.style.display = "none";
		    contour_video_btn.removeAttribute("disabled");

		    feedback_el.innerText = "Failed to contour still, " + err;		    
		});
		
	    }, 10);
	    
	    return false;
	};
	
	image_btn.onclick = function(){

	    feedback_el.innerHTML = "";
	    
	    try {
		process_upload();
	    } catch(err) {
		feedback_el.innerText = "Failed to process upload, " + err;
		console.error(err);
	    }

	    return false;
	};

	start_video_btn.onclick = function(){

	    feedback_el.innerHTML = "";
	    
	    navigator.mediaDevices.getUserMedia({ video: { facingMode: "environment" } }).then(function(stream) {

		stop_video_btn.onclick = function(){

		    video.pause();
		    video.srcObject = null;
		    
		    stream.getTracks().forEach((track) => {
			if (track.readyState == 'live') {
			    track.stop();
			}
		    });

		    stop_video_btn.setAttribute("disabled", "disabled");
		    start_video_btn.removeAttribute("disabled");
		    contour_video_btn.removeAttribute("disabled");		    		    
		};

		start_video_btn.setAttribute("disabled", "disabled");
		contour_video_btn.removeAttribute("disabled");				
		stop_video_btn.removeAttribute("disabled");
		
		process_video(stream);
		
	    }).catch((err) => {
		feedback_el.innerText = "Failed to start video feed, " + err;
	    });
	};
	
	upload_el.removeAttribute("disabled");
	image_btn.removeAttribute("disabled");
	start_video_btn.removeAttribute("disabled");	
	
    }).catch((err) => {
	feedback_el.innerText = "Failed to load age WebAssembly functions, " + err;
        return false;
    });
	
});
