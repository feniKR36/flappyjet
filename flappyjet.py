import pygame
import random
import os
pygame.init()
#model
screenwidth = 400
screenheight = 600
jetsize = 20
towerwidth = 50
towergap = 200
jetspeed = 5
gravity = 0.5
jetstrength = -10
fps = 60
scorefile = "scores.txt" 
#insert towers and jet
WHITE = (255,255,255)
BLACK = (0,0,0)
GREEN = (0,200,0)
screen = pygame.display.set_mode((screenwidth,screenheight))
pygame.display.set_caption("Ready...Jet set go!")
clock = pygame.time.Clock()
font = pygame.font.Font(None,36)
#scorerecs
def load_scores():
    scores = {}
    if not os.path.exists(scorefile):
        return scores
    with open(scorefile,"r") as f:
        for line in f:
            name,score = line.strip().split(",")
            scores[name] = int(score)
    return scores
def save_scores(scores):
    with open(scorefile,"w") as f:
        for name,score in scores.items():
            f.write(f"{name},{score}\n")
def update_player_score(username,new_score):
    scores = load_scores()
    if username not in scores or new_score > scores[username]:
        scores[username] = new_score
    save_scores(scores)
def get_player_best(username):
    scores = load_scores()
    return scores.get(username,0)
def get_leaderboard():
    scores = load_scores()
    return sorted(scores.items(),key=lambda x:x[1],reverse=True)
#logindeets
def login_screen():
    username = ""
    while True:
        screen.fill(WHITE)
        title = font.render("Enter Username",True,BLACK)
        name = font.render(username,True,BLACK)
        screen.blit(title,(120,250))
        screen.blit(name,(120,300))
        pygame.display.flip()
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                pygame.quit()
                exit()
            if event.type == pygame.KEYDOWN:
                if event.key == pygame.K_RETURN and username.strip():
                    return username.strip()
                elif event.key == pygame.K_BACKSPACE:
                    username = username[:-1]
                else:
                    if len(username) < 12:
                        username += event.unicode
#mainwindowmenu
def main_menu(username):
    options = ["Play","Leaderboard","Profile","Quit"]
    selected = 0
    while True:
        screen.fill(WHITE)
        title = font.render(f"Welcome {username}",True,BLACK)
        screen.blit(title,(120,120))
        for i,option in enumerate(options):
            color = GREEN if i == selected else BLACK
            text = font.render(option,True,color)
            screen.blit(text,(150,220 + i*50))
        pygame.display.flip()
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                pygame.quit()
                exit()
            if event.type == pygame.KEYDOWN:
                if event.key == pygame.K_UP:
                    selected = (selected - 1) % len(options)
                elif event.key == pygame.K_DOWN:
                    selected = (selected + 1) % len(options)
                elif event.key == pygame.K_RETURN:
                    return options[selected]
#profiledeets
def profile_screen(username):
    best_score = get_player_best(username)
    while True:
        screen.fill(WHITE)
        title = font.render("Player Profile",True,BLACK)
        name = font.render(f"Username: {username}",True,BLACK)
        best = font.render(f"Best Score: {best_score}",True,BLACK)
        back = font.render("ESC to return",True,BLACK)
        screen.blit(title,(120,150))
        screen.blit(name,(120,230))
        screen.blit(best,(120,270))
        screen.blit(back,(120,420))
        pygame.display.flip()
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                pygame.quit()
                exit()
            if event.type == pygame.KEYDOWN:
                if event.key == pygame.K_ESCAPE:
                    return

def scoreboard_screen():#scorerecsdeets
    leaderboard = get_leaderboard()
    while True:
        screen.fill(WHITE)
        title = font.render("High Scores",True,BLACK)
        screen.blit(title,(130,80))
        for i,(name,score) in enumerate(leaderboard[:5]):
            text = font.render(f"{i+1}. {name} : {score}",True,BLACK)
            screen.blit(text,(110,150 + i*40))
        play_text = font.render("SPACE = menu",True,BLACK)
        reset_text = font.render("R = reset scores",True,BLACK)
        screen.blit(play_text,(100,500))
        screen.blit(reset_text,(90,540))
        pygame.display.flip()
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                pygame.quit()
                exit()
            if event.type == pygame.KEYDOWN:
                if event.key == pygame.K_SPACE:
                    return
                if event.key == pygame.K_r:
                    save_scores({})
                    leaderboard = []

def create_tower():
    height = random.randint(50,screenheight-towergap-50)
    return {
        "x":screenwidth,
        "height":height,
        "scored":False
    }

def run_game(username):#looping
    jet_yaxis = screenheight//2
    jet_velocity = 0
    towers = []
    score = 0
    best_score = get_player_best(username)
    
    while True:
        screen.fill(WHITE)
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                pygame.quit()
                exit()
            if event.type == pygame.KEYDOWN or event.type == pygame.MOUSEBUTTONDOWN:
                jet_velocity = jetstrength
        jet_velocity += gravity #jetphysics
        jet_yaxis += jet_velocity
        jet_rect = pygame.Rect(50,jet_yaxis,jetsize,jetsize)
        pygame.draw.rect(screen,BLACK,jet_rect)
        if not towers or towers[-1]["x"] < screenwidth-200:
            towers.append(create_tower())
        for tower in towers:
            tower["x"] -= jetspeed
            top_tower = pygame.Rect(tower["x"],0,towerwidth,tower["height"])
            bottom_tower = pygame.Rect(tower["x"],tower["height"]+towergap,towerwidth,screenheight)
            pygame.draw.rect(screen,GREEN,top_tower)
            pygame.draw.rect(screen,GREEN,bottom_tower)
            if jet_rect.colliderect(top_tower) or jet_rect.colliderect(bottom_tower):
                update_player_score(username,score)
                return
            if tower["x"] + towerwidth < 50 and not tower["scored"]:
                score += 1
                tower["scored"] = True
        if towers and towers[0]["x"] < -towerwidth:
            towers.pop(0)
        if jet_yaxis < 0 or jet_yaxis > screenheight:
            update_player_score(username,score)
            return
        score_text = font.render(f"Score: {score}",True,BLACK)
        best_text = font.render(f"Best: {best_score}",True,BLACK)
        screen.blit(score_text,(10,10))
        screen.blit(best_text,(280,10))
        pygame.display.flip()
        clock.tick(fps)
username = login_screen()#mainprog
while True:
    choice = main_menu(username)
    if choice == "Play":
        run_game(username)
    elif choice == "Leaderboard":
        scoreboard_screen()
    elif choice == "Profile":
        profile_screen(username)
    elif choice == "Quit":
        pygame.quit()
        break