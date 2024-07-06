package service

import "fmt"

// SayHello 实现了 GreetingService 的 SayHello 方法
func (s *SimpleGreetingService) SayHello(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

func (s *SimpleGreetingService) GenerateSummary(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

func (s *SimpleGreetingService) ImportFileSummary(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}
