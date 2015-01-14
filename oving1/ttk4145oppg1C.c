#include <pthread.h>
#include <stdio.h>

int i = 0;

void* function1(){
	for (int j=0; j<1000000; j++){
		i++;
	}
	return NULL;
}

void* function2(){
	for (int j=0; j<1000000; j++){
		i--;
	}
	return NULL;
}

int main(){
	
	pthread_t thread1;
	pthread_create(&thread1, NULL, function1, NULL);

	pthread_t thread2;
	pthread_create(&thread2, NULL, function2, NULL);

	pthread_join(thread1, NULL);
	pthread_join(thread2, NULL);

	printf("%i\n", i);

	return 0;
}
