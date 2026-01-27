import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Service } from './entities/service.entity';
import { CreateServiceDto } from './dto/create-service.dto';
import { UpdateServiceDto } from './dto/update-service.dto';

@Injectable()
export class ServicesService {
  constructor(@InjectRepository(Service) private repo: Repository<Service>) {}

  async create(dto: CreateServiceDto): Promise<Service> {
    const service = this.repo.create(dto);
    return await this.repo.save(service);
  }

  async findAll(): Promise<any> {
    const services = await this.repo.find({ where: { published: true }, order: { createdAt: 'DESC' } });
    return { success: true, data: services };
  }

  async findBySlug(slug: string): Promise<any> {
    const service = await this.repo.findOne({ where: { slug, published: true } });
    if (!service) throw new NotFoundException(`Service "${slug}" not found`);
    return { success: true, data: service };
  }

  async update(id: string, dto: UpdateServiceDto): Promise<Service> {
    const service = await this.repo.findOne({ where: { id } });
    if (!service) throw new NotFoundException(`Service "${id}" not found`);
    Object.assign(service, dto);
    return await this.repo.save(service);
  }

  async remove(id: string): Promise<void> {
    const service = await this.repo.findOne({ where: { id } });
    if (!service) throw new NotFoundException(`Service "${id}" not found`);
    await this.repo.remove(service);
  }
}
