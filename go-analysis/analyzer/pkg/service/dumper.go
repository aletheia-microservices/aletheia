package service

import "analyzer/pkg/logger"

func (node *Service) Yaml() map[string]interface{} {
	data := make(map[string]interface{})
	// exposed methods
	for k, method := range node.ExposedMethods {
		if method == nil {
			logger.Logger.Fatalf("[YAML] [%s] unexpected nil exposed method for name (%s) in list:\n%v", node.Name, k, node.ExposedMethods)
		}
		data[method.String()], _ = method.Yaml()
	}
	// queue handler methods
	for k, method := range node.QueueHandlerMethods {
		if method == nil {
			logger.Logger.Fatalf("[YAML] [%s] unexpected nil queue handler method for name (%s) in list:\n%v", node.Name, k, node.QueueHandlerMethods)
		}
		data[method.String()], _ = method.Yaml()
	}
	// internal methods (already contain handlers)
	for k, method := range node.InternalMethods {
		if method == nil {
			logger.Logger.Fatalf("[YAML] [%s] unexpected nil internal method for name (%s) in list:\n%v", node.Name, k, node.InternalMethods)
		}
		data[method.String()], _ = method.Yaml()
	}
	for k, method := range node.PackageMethods {
		if method == nil {
			logger.Logger.Fatalf("[YAML] [%s] unexpected nil package method for name (%s) in list:\n%v", node.Name, k, node.PackageMethods)
		}
		data[method.String()], _ = method.Yaml()
	}
	return data
}

func (node *Service) YamlCalls() map[string][]string {
	data := make(map[string][]string)
	// exposed methods
	for k, method := range node.ExposedMethods {
		if method == nil {
			logger.Logger.Fatalf("[YAML] [%s] unexpected nil exposed method for name (%s) in list:\n%v", node.Name, k, node.ExposedMethods)
		}
		data[method.String()] = method.YamlCalls()
	}
	// queue handler methods
	for k, method := range node.QueueHandlerMethods {
		if method == nil {
			logger.Logger.Fatalf("[YAML] [%s] unexpected nil queue handler method for name (%s) in list:\n%v", node.Name, k, node.QueueHandlerMethods)
		}
		data[method.String()] = method.YamlCalls()
	}
	// internal methods (already contain handlers)
	for k, method := range node.InternalMethods {
		if method == nil {
			logger.Logger.Fatalf("[YAML] [%s] unexpected nil internal method for name (%s) in list:\n%v", node.Name, k, node.InternalMethods)
		}
		data[method.String()] = method.YamlCalls()
	}
	for k, method := range node.PackageMethods {
		if method == nil {
			logger.Logger.Fatalf("[YAML] [%s] unexpected nil package method for name (%s) in list:\n%v", node.Name, k, node.PackageMethods)
		}
		data[method.String()] = method.YamlCalls()
	}
	return data
}
