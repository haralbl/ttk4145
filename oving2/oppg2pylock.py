#from threading import Thread
import threading

i = 0
lock = threading.Lock()

def function1():
	global i
	with lock:
		for j in range(1000000):
			i+=1
		
def function2():
	global i
	with lock:
		for j in range(1000001):
			i-=1

def main():
	global i
	thread1 = threading.Thread(target = function1, args = (),)
	thread1.start()
	
	thread2 = threading.Thread(target = function2, args = (),)
	thread2.start()
	
	thread1.join()
	thread2.join()
	
	print(i)

main()
