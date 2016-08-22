(function() {
	'use strict';
	var c = document.getElementById("chess_panel");
	var ctx = c.getContext("2d");

	var $c = $(c);
	var width = $c.attr("width")
	var height = $c.attr("height")

	var logDiv = document.getElementById("log")


	var margin = 15;

	var PLAYER_1 = 1;
	var PLAYER_2 = 2;
	var EMPTY = undefined;

	var x_max = 19;
	var y_max = 19;

	var player = 1;

	var panel = new Array(19);
	for (var i = 0; i < 19; i++) {
		panel[i] = new Array(19);
	}

	var game_logs = [];


	var is_game_started = false;
	var is_player_turn = false;
	var is_ai_first = false;
	var player_name = "";

	var horizontal_space = ((height - margin) - (10 + margin)) / (19 - 1);
	var vertical_space = ((width - margin) - (10 + margin)) / (19 - 1);
	console.log("horizontal_space: " + horizontal_space);


	drawHorizontalLine(ctx, 10 + margin, 10 + margin, width - margin, height - margin, 18)
	drawVerticalLine(ctx, 10 + margin, 10 + margin, width - margin, height - margin, 18)
	drawColumnName(ctx, 10 - 0.5 + margin, margin, 18)
	drawRowName(ctx, 0, 10 + margin, 18)


	c.addEventListener('click', function(event) {
		if (!is_game_started) {
			$("#status_display").html("遊戲尚未開始！");
			return;
		}
		if (!is_player_turn) {
			$("#status_display").html("還不是玩家的回合！");
			return;

		}
		var real_x = event.pageX - c.offsetLeft;
		var real_y = event.pageY - c.offsetTop;
		console.log("real: (" + real_x + ", " + real_y + ")")

		var x = Math.floor((real_x - (10 + margin)) / vertical_space + 0.5);
		var y = Math.floor((real_y - (10 + margin)) / horizontal_space + 0.5);
		console.log("Set: (" + x + ", " + y + ")")


		var success = set(x, y, is_ai_first ? PLAYER_2 : PLAYER_1);
		if (success != true) {
			console.log("ERR");
			return;
		}
		is_player_turn = false;
		if (is_game_started) {
			callAi();
		}

	});

	document.getElementById("single_player").addEventListener("click", function(event) {
		$("#choose_player_number").hide();
		$("#input_player_name").show();
	})

	document.getElementById("player_name_submit").addEventListener("click", function(event) {
		$("#input_player_name").hide();
		$("#choose_who_first").show();
		$("#player_name_display").html(player_name = $("#player_name").val());
	})

	document.getElementById("ai_first").addEventListener("click", function(event) {
		is_ai_first = true;
		whoFirstChossen(event);
		is_player_turn = false;
		callAi();
	})

	document.getElementById("player_first").addEventListener("click", function(event) {
		is_ai_first = false;
		whoFirstChossen(event);
		is_player_turn = true;
	})

	function whoFirstChossen(event) {
		$("#choose_who_first").hide();
		$("#game_status").show();
		is_game_started = true;
		$("#status_display").html("遊戲開始");
	}


	function drawHorizontalLine(ctx, startx, starty, endx, endy, line_n) {
		for (var i = 0; i <= line_n; i++) {
			ctx.moveTo(startx, starty + i * horizontal_space);
			ctx.lineTo(endx, starty + i * horizontal_space);
			ctx.stroke();

		}

	}

	function drawVerticalLine(ctx, startx, starty, endx, endy, line_n) {
		for (var i = 0; i <= line_n; i++) {
			ctx.moveTo(startx + i * vertical_space, starty);
			ctx.lineTo(startx + i * vertical_space, endy);
			ctx.stroke();

		}
	}

	function drawColumnName(ctx, startx, starty, line_n) {
		for (var i = 0; i <= line_n; i++) {
			ctx.font = "9px Arial";
			ctx.fillText(i + 1, startx + i * vertical_space, starty);
		}
	}

	function drawRowName(ctx, startx, starty, line_n) {
		for (var i = 0; i <= line_n; i++) {
			ctx.font = "9px Arial";
			ctx.fillText(i + 1, startx, starty + i * horizontal_space);
		}
	}

	function drawStone(ctx, x, y, color) {
		ctx.beginPath();
		var rad = Math.min(horizontal_space, vertical_space) / 2
		rad = rad * 0.6;
		ctx.arc(x, y, rad, 0, 2 * Math.PI);
		ctx.fillStyle = color;
		ctx.fill();
		ctx.stroke();
	}

	function callAi() {
		console.log(game_logs);
		$.ajax({
			// url: "/random_ai",
			url: "/pichu_ai",
			type: 'post',
			dataType: 'json',
			data: JSON.stringify(game_logs),
			success: function(res) {
				console.log(res);
				set(res[0], res[1], is_ai_first ? PLAYER_1 : PLAYER_2);
				is_player_turn = true;

			}


		});
	}

	function set(x, y, player) {
		if (panel[x][y] != EMPTY) {
			console.log("壓到了" + panel[x][y]);
			$("#status_display").html("壓到了!!");
			return false
		}

		logDiv.innerHTML = logDiv.innerHTML + "Player " + player + ": Put (" + (x + 1) + ", " + (y + 1) + ")<br/>";

		panel[x][y] = player;
		game_logs.push([x, y]);
		drawStone(ctx, (10 + margin) + (x) * vertical_space, (10 + margin) + (y) * horizontal_space, player == PLAYER_2 ? "white" : "black")
		player = player == 1 ? 2 : 1

		var status = judge(x, y);

		if (status != 0) {
			if (
				(status == PLAYER_1 && is_ai_first) || (status == PLAYER_2 && !is_ai_first)
			) {
				$("#status_display").html("Pichu's AI 獲勝");
			} else if (
				(status == PLAYER_1 && !is_ai_first) || (status == PLAYER_2 && is_ai_first)
			) {
				$("#status_display").html("玩家 " + player_name + "獲勝");
			} else if (status == 3) {
				$("#status_display").html("平手");
			}
			is_game_started = false;
		}
		return true;
	}

	function judge(x, y) {
		var win = 0;
		if (panel[x][y] == EMPTY) {
			return 0; // no result
		}

		var k, num;
		num = 0;
		k = 1;
		var b = panel;

		// Right
		while ((x + k < x_max) && (b[x][y] == b[x + k][y])) {
			if (++num >= 4) {
				console.log("Found Winner");
				return b[x][y];
			}
			k++;
		}
		k = 1;
		// Left
		while ((x - k >= 0) && (b[x][y] == b[x - k][y])) {
			if (++num >= 4) {
				console.log("Found Winner");
				return b[x][y];
			}
			k++;
		}

		num = 0;
		k = 1;
		// Down
		while ((y + k < y_max) && (b[x][y] == b[x][y + k])) {
			if (++num >= 4) {
				console.log("Found Winner");
				return panel[x][y];
			}
			k++;
		}
		k = 1;
		// UP
		while ((y - k >= 0) && (b[x][y] == b[x][y - k])) {
			if (++num >= 4) {
				console.log("Found Winner");
				return panel[x][y];
			}
			k++;
		} //2    
		num = 0;
		k = 1;

		// Right and Down
		while ((x + k < x_max) && (y + k < y_max) && (b[x][y] == b[x + k][y + k])) {
			if (++num >= 4) {
				console.log("Found Winner");
				return panel[x][y];
			}
			k++;
		}
		k = 1;
		// Left and Up
		while ((x - k >= 0) && (y - k >= 0) && (b[x][y] == b[x - k][y - k])) {
			if (++num >= 4) {
				console.log("Found Winner");
				return panel[x][y];
			}
			k++;
		} //3
		num = 0;
		k = 1;

		// Right and Up
		while ((x + k < x_max) && (y - k >= 0) && (b[x][y] == b[x + k][y - k])) {
			if (++num >= 4) {
				console.log("Found Winner");
				return panel[x][y];
			}
			k++;
		}
		k = 1;
		// Left and Down
		while ((x - k >= 0) && (y + k < y_max) && (b[x][y] == b[x - k][y + k])) {
			if (++num >= 4) {
				console.log("Found Winner");
				return panel[x][y];
			}
			k++;
		}

		// Check panel full
		for (var cy = 0; cy < y_max; cy++) {
			for (var cx = 0; cx < x_max; cx++) {
				if (b[cx][cy] == EMPTY) {
					return 0;
				}
			}
		}
		return 3;

	}


})();
