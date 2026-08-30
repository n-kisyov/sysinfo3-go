export namespace collector {
	
	export class BIOSInfo {
	    vendor: string;
	    version: string;
	    date: string;
	    manufacturer: string;
	    model: string;
	    serial_number?: string;
	
	    static createFrom(source: any = {}) {
	        return new BIOSInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vendor = source["vendor"];
	        this.version = source["version"];
	        this.date = source["date"];
	        this.manufacturer = source["manufacturer"];
	        this.model = source["model"];
	        this.serial_number = source["serial_number"];
	    }
	}
	export class BatteryInfo {
	    percentage: number;
	    status: string;
	    time_left?: string;
	
	    static createFrom(source: any = {}) {
	        return new BatteryInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.percentage = source["percentage"];
	        this.status = source["status"];
	        this.time_left = source["time_left"];
	    }
	}
	export class CPUInfo {
	    model: string;
	    physical_cores: number;
	    logical_cores: number;
	    usage_percent: number;
	    per_core_usage: number[];
	    temperature: string;
	
	    static createFrom(source: any = {}) {
	        return new CPUInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.model = source["model"];
	        this.physical_cores = source["physical_cores"];
	        this.logical_cores = source["logical_cores"];
	        this.usage_percent = source["usage_percent"];
	        this.per_core_usage = source["per_core_usage"];
	        this.temperature = source["temperature"];
	    }
	}
	export class DiskInfo {
	    drive_letter: string;
	    mount_point: string;
	    fs_type: string;
	    total_bytes: number;
	    total: string;
	    used_bytes: number;
	    used: string;
	    free_bytes: number;
	    free: string;
	    used_percent: number;
	
	    static createFrom(source: any = {}) {
	        return new DiskInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.drive_letter = source["drive_letter"];
	        this.mount_point = source["mount_point"];
	        this.fs_type = source["fs_type"];
	        this.total_bytes = source["total_bytes"];
	        this.total = source["total"];
	        this.used_bytes = source["used_bytes"];
	        this.used = source["used"];
	        this.free_bytes = source["free_bytes"];
	        this.free = source["free"];
	        this.used_percent = source["used_percent"];
	    }
	}
	export class GPUInfo {
	    name: string;
	    driver: string;
	    vram_bytes: number;
	    vram: string;
	
	    static createFrom(source: any = {}) {
	        return new GPUInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.driver = source["driver"];
	        this.vram_bytes = source["vram_bytes"];
	        this.vram = source["vram"];
	    }
	}
	export class HostInfo {
	    hostname: string;
	    os: string;
	    platform: string;
	    platform_version: string;
	    kernel_version: string;
	    uptime_seconds: number;
	    uptime: string;
	
	    static createFrom(source: any = {}) {
	        return new HostInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hostname = source["hostname"];
	        this.os = source["os"];
	        this.platform = source["platform"];
	        this.platform_version = source["platform_version"];
	        this.kernel_version = source["kernel_version"];
	        this.uptime_seconds = source["uptime_seconds"];
	        this.uptime = source["uptime"];
	    }
	}
	export class MemoryInfo {
	    total_bytes: number;
	    total: string;
	    used_bytes: number;
	    used: string;
	    available_bytes: number;
	    available: string;
	    used_percent: number;
	    swap_total_bytes: number;
	    swap_total: string;
	    swap_used_bytes: number;
	    swap_used: string;
	    swap_used_percent: number;
	
	    static createFrom(source: any = {}) {
	        return new MemoryInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_bytes = source["total_bytes"];
	        this.total = source["total"];
	        this.used_bytes = source["used_bytes"];
	        this.used = source["used"];
	        this.available_bytes = source["available_bytes"];
	        this.available = source["available"];
	        this.used_percent = source["used_percent"];
	        this.swap_total_bytes = source["swap_total_bytes"];
	        this.swap_total = source["swap_total"];
	        this.swap_used_bytes = source["swap_used_bytes"];
	        this.swap_used = source["swap_used"];
	        this.swap_used_percent = source["swap_used_percent"];
	    }
	}
	export class NetInterface {
	    name: string;
	    mac: string;
	    addresses: string[];
	    mtu: number;
	    flags: string[];
	    bytes_sent: number;
	    bytes_recv: number;
	    bytes_sent_h: string;
	    bytes_recv_h: string;
	    bytes_sent_per_sec: number;
	    bytes_recv_per_sec: number;
	    bytes_sent_per_sec_h: string;
	    bytes_recv_per_sec_h: string;
	
	    static createFrom(source: any = {}) {
	        return new NetInterface(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.mac = source["mac"];
	        this.addresses = source["addresses"];
	        this.mtu = source["mtu"];
	        this.flags = source["flags"];
	        this.bytes_sent = source["bytes_sent"];
	        this.bytes_recv = source["bytes_recv"];
	        this.bytes_sent_h = source["bytes_sent_h"];
	        this.bytes_recv_h = source["bytes_recv_h"];
	        this.bytes_sent_per_sec = source["bytes_sent_per_sec"];
	        this.bytes_recv_per_sec = source["bytes_recv_per_sec"];
	        this.bytes_sent_per_sec_h = source["bytes_sent_per_sec_h"];
	        this.bytes_recv_per_sec_h = source["bytes_recv_per_sec_h"];
	    }
	}
	export class PhysicalDiskInfo {
	    name: string;
	    model: string;
	    type: string;
	    size_bytes: number;
	    size: string;
	    read_bytes: number;
	    write_bytes: number;
	    read_bytes_h: string;
	    write_bytes_h: string;
	    read_bytes_per_sec: number;
	    write_bytes_per_sec: number;
	    read_bytes_per_sec_h: string;
	    write_bytes_per_sec_h: string;
	
	    static createFrom(source: any = {}) {
	        return new PhysicalDiskInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.model = source["model"];
	        this.type = source["type"];
	        this.size_bytes = source["size_bytes"];
	        this.size = source["size"];
	        this.read_bytes = source["read_bytes"];
	        this.write_bytes = source["write_bytes"];
	        this.read_bytes_h = source["read_bytes_h"];
	        this.write_bytes_h = source["write_bytes_h"];
	        this.read_bytes_per_sec = source["read_bytes_per_sec"];
	        this.write_bytes_per_sec = source["write_bytes_per_sec"];
	        this.read_bytes_per_sec_h = source["read_bytes_per_sec_h"];
	        this.write_bytes_per_sec_h = source["write_bytes_per_sec_h"];
	    }
	}
	export class ProcessInfo {
	    pid: number;
	    name: string;
	    cpu_percent: number;
	    memory_bytes: number;
	    memory: string;
	
	    static createFrom(source: any = {}) {
	        return new ProcessInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pid = source["pid"];
	        this.name = source["name"];
	        this.cpu_percent = source["cpu_percent"];
	        this.memory_bytes = source["memory_bytes"];
	        this.memory = source["memory"];
	    }
	}
	export class SystemSnapshot {
	    // Go type: time
	    timestamp: any;
	    host: HostInfo;
	    memory: MemoryInfo;
	    cpu: CPUInfo;
	    disks: DiskInfo[];
	    physical_disks: PhysicalDiskInfo[];
	    network: NetInterface[];
	    battery?: BatteryInfo;
	    processes: ProcessInfo[];
	    gpu: GPUInfo[];
	    bios: BIOSInfo;
	
	    static createFrom(source: any = {}) {
	        return new SystemSnapshot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.host = this.convertValues(source["host"], HostInfo);
	        this.memory = this.convertValues(source["memory"], MemoryInfo);
	        this.cpu = this.convertValues(source["cpu"], CPUInfo);
	        this.disks = this.convertValues(source["disks"], DiskInfo);
	        this.physical_disks = this.convertValues(source["physical_disks"], PhysicalDiskInfo);
	        this.network = this.convertValues(source["network"], NetInterface);
	        this.battery = this.convertValues(source["battery"], BatteryInfo);
	        this.processes = this.convertValues(source["processes"], ProcessInfo);
	        this.gpu = this.convertValues(source["gpu"], GPUInfo);
	        this.bios = this.convertValues(source["bios"], BIOSInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

